package pg

import (
	"miniDouyin/biz/dal/rdb"
	"miniDouyin/biz/model/miniDouyin/api"
	"miniDouyin/utils"
	"reflect"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64 `gorm:"primaryKey"`
	VideoId   int64
	UserId    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt `gorm:"default:NULL"`
}

func (v *Comment) TableName() string {
	return "comments"
}

// 数据库模型转换为api的结构体
func (v *Comment) ToApiComment(cUser *DBUser, clientUser *DBUser) (*api.Comment, error) {
	ac := &api.Comment{
		ID:         v.ID,
		User:       nil,
		Content:    v.Content,
		CreateDate: v.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 填充评论用户信息
	capiUser, err := cUser.ToApiUser(clientUser)
	if err != nil {
		return ac, err
	}
	ac.User = capiUser

	return ac, nil
}

// 发布评论
func (v *Comment) CreateComment() (int64, error) {
	result := DB.Create(v)
	return v.ID, result.Error
}

// 更新缓存
func (u *Comment) UpdateRedis(Type int) {
	items := utils.StructToMap(u)
	msg := RedisMsg{
		TYPE: Type,
		DATA: items,
	}
	ChanFromDB <- msg

}

// 删除评论
func DeleteComment(commentId int64) (*Comment, error) {
	comment := Comment{}
	// 先查询是否有此评论
	result := DB.Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return nil, utils.ErrDelCommentNotExist
	}

	// 如果评论存在则删除此条评论
	result = DB.Where("id = ?", commentId).Delete(&Comment{})
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

// 从获取用户信息请求请求构造新用户
func DBGetCommentByID(CommentID int64) (*Comment, error) {
	var comm Comment
	res := DB.First(&comm, "ID = ?", CommentID)

	if res.Error != nil {
		// 没有找到记录
		return nil, utils.ErrUserNotFound
	}

	return &comm, nil
}

// 根据 video_id获取评论列表
func GetDBCommentList(v_id int64) (clist []Comment, err error) {
	err = nil
	res := DB.Model(&Comment{}).Where("video_id = ?", v_id).
		Order("ID").Find(&clist)
	if res.Error != nil {
		err = utils.ErrGetCommentListFailed
	}
	return
}

func ID2Comment(cID int64) (*Comment, error) {
	// 尝试从缓存查询评论
	var cDB = &Comment{}
	var err error
	cMap, find := rdb.GetCommentByID(strconv.FormatInt(cID, 10))
	if find {
		// 缓存命中
		log.Debugln("ID2Comment: 从缓存查询评论成功")
		cDB = &Comment{}
		cDB.InitSelfFromMap(cMap)
	} else {
		//从数据库查找
		cDB, err = DBGetCommentByID(cID)
		if err != nil {
			return nil, utils.ErrCommentNotExist
		}
		// 发送消息更新缓存
		cDB.UpdateRedis(CommentCreate)
		log.Infoln("ID2Comment：更新comment缓存")
	}
	return cDB, nil
}

// 获取评论列表
func NewGetCommentListService(v_id int64, token string) (clist []*api.Comment, r_err error) {
	clientUser, err := Token2DBUser(token)
	if err != nil {
		return nil, err
	}

	// 校验视频id合法性
	v := &DBVideo{}
	res := DB.Model(v).First(v, "ID = ?", v_id)
	if res.Error != nil {
		return nil, utils.ErrVideoNotExist
	}

	// 如果视频id有效再获取评论列表
	// 缓存未命中，从数据库查
	commentlist, find := rdb.GetVideoCommentList(int(v_id))
	if find {
		// 缓存命中
		log.Debugln("GetCommentList: 从缓存查询评论列表成功")
		for _, c_id := range commentlist {
			// 尝试从查询评论
			C_ID, _ := strconv.Atoi(c_id)
			cDB, err := ID2Comment(int64(C_ID))
			if err != nil {
				return nil, err
			}

			// 查询评论者
			cUser, err := ID2DBUser(cDB.UserId)
			if err != nil {
				return nil, err
			}
			ac, err := cDB.ToApiComment(cUser, clientUser)
			if err != nil {
				return nil, utils.ErrGetCommentListFailed
			}
			clist = append(clist, ac)
		}
		return
	} else {
		// 从数据库查
		cDBlist, err := GetDBCommentList(v_id)
		if err != nil {
			return nil, err
		}
		// 将评论列表格式进行转换
		for _, dbcomment := range cDBlist {
			// 更新缓存
			dbcomment.UpdateRedis(CommentCreate)

			cUser, err := ID2DBUser(dbcomment.UserId)
			if err != nil {
				return nil, utils.ErrUserNotFound
			}
			ac, err := dbcomment.ToApiComment(cUser, clientUser)
			if err != nil {
				return nil, utils.ErrGetCommentListFailed
			}
			clist = append(clist, ac)
		}
		return
	}
}

// 根据评论请求封装发布评论
func dbCreateComment(req *api.CommentActionRequest, userId int64) (*Comment, error) {
	comm := &Comment{
		VideoId: req.VideoID,
		UserId:  userId,
		Content: *req.CommentText,
	}
	_, err := comm.CreateComment()
	if err != nil {
		return nil, err
	}
	return comm, nil
}

func (u *Comment) InitSelfFromMap(uMap map[string]string) {
	reflectVal := reflect.ValueOf(u).Elem()

	for fieldName, fieldValue := range uMap {
		field := reflectVal.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.Int, reflect.Int64:
				tmp, _ := strconv.ParseInt(fieldValue, 10, 64)
				field.SetInt(tmp)
			case reflect.String:
				field.SetString(fieldValue)
			}
		}
	}
}
