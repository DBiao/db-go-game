package model

import "db-go-game/pkg/entity"

type User struct {
	entity.GormEntityTs
	Uid         int64  `gorm:"column:uid;primary_key" json:"uid"`                          // 用户ID 系统生成
	Account     string `gorm:"column:lark_id;uniqueIndex;NOT NULL" json:"account"`         // 账户ID 用户设置
	Password    string `gorm:"column:password;NOT NULL" json:"password"`                   // 密码
	Uuid        string `gorm:"column:uuid;NOT NULL" json:"uuid"`                           // 注册设备唯一标识
	Status      int    `gorm:"column:status;default:0;NOT NULL" json:"status"`             // 用户状态
	Nickname    string `gorm:"column:nickname;NOT NULL" json:"nickname"`                   // 昵称
	Gender      int    `gorm:"column:gender;default:0;NOT NULL" json:"gender"`             // 性别
	BirthTs     int64  `gorm:"column:birth_ts;default:0;NOT NULL" json:"birth_ts"`         // 生日
	Email       string `gorm:"column:email;NOT NULL" json:"email"`                         // Email
	Mobile      string `gorm:"column:mobile;NOT NULL" json:"mobile"`                       // 手机号
	RegPlatform int    `gorm:"column:reg_platform;default:0;NOT NULL" json:"reg_platform"` // 注册平台
	CityId      int    `gorm:"column:city_id;default:0;NOT NULL" json:"city_id"`           // 城市ID
	AvatarKey   string `gorm:"column:avatar_key" json:"avatar_key"`                        // 小图 72*72
}

func (*User) TableName() string {
	return "user"
}
