package service

import (
	"db-go-game/domain/dao"
	"db-go-game/domain/model"
	"db-go-game/pkg/common/djwt"
	"db-go-game/pkg/common/dlog"
	"db-go-game/pkg/common/dredis"
	"db-go-game/pkg/common/dsnowflake"
	"db-go-game/pkg/constant"
	"db-go-game/pkg/dhttp"
	"db-go-game/pkg/entity"
	"db-go-game/pkg/utils"
	"db-go-game/services/api/client"
	"db-go-game/services/api/internal/config"
	"db-go-game/services/api/internal/dto"
	"github.com/jinzhu/copier"
)

type IAuthService interface {
	SignUp(req *dto.SignUpReq) *dhttp.Resp
	SignIn(req *dto.SignInReq) *dhttp.Resp
	RefreshToken(token string) *dhttp.Resp
	SignOut(token string) *dhttp.Resp
}

type AuthService struct {
	userDao     dao.IUserDao
	logicClient client.ILogicClient
}

func NewAuthService(AuthDao dao.IUserDao) IAuthService {
	conf := config.GetConfig()
	return &AuthService{
		userDao:     AuthDao,
		logicClient: client.NewLogicClient(conf.Etcd, conf.LogicServer, conf.Jaeger, conf.Name),
	}
}

func (u *AuthService) SignUp(req *dto.SignUpReq) *dhttp.Resp {
	resp := new(dhttp.Resp)

	Auth := &model.User{}
	copier.Copy(Auth, req)
	Auth.Uid = dsnowflake.NewSnowflakeID()
	Auth.Password = utils.MD5(req.Password)

	if err := u.userDao.Create(Auth); err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	return resp
}

func (u *AuthService) SignIn(req *dto.SignInReq) *dhttp.Resp {
	var (
		resp  = new(dhttp.Resp)
		w     = entity.NewMysqlWhere()
		err   error
		token *djwt.JwtToken
	)

	w.SetFilter("account = ?", req.Account)
	w.SetFilter("password = ?", utils.MD5(req.Password))

	Auth, err := u.userDao.VerifyIdentity(w)
	if err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	if Auth.Uid == 0 {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	token, err = djwt.CreateToken(Auth.Uid, true, constant.CONST_DURATION_SHA_JWT_ACCESS_TOKEN_EXPIRE_IN_SECOND)
	if err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	resp.Data = token
	return resp
}

func (u *AuthService) RefreshToken(token string) *dhttp.Resp {
	resp := new(dhttp.Resp)

	jwtToken, err := djwt.Decode(token)
	if err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	//if time.Now().Unix() > jwtToken.Expire {
	//	resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
	//	dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
	//	return resp
	//}

	newToken, err := djwt.CreateToken(jwtToken.Uid, true, constant.CONST_DURATION_SHA_JWT_ACCESS_TOKEN_EXPIRE_IN_SECOND)
	if err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	resp.Data = newToken
	return resp
}

func (u *AuthService) SignOut(token string) *dhttp.Resp {
	resp := new(dhttp.Resp)

	jwtToken, err := djwt.Decode(token)
	if err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_REGISTER_FAILED, err.Error())
		return resp
	}

	key := constant.RK_SYNC_USER_ACCESS_TOKEN + utils.Int64ToStr(jwtToken.Uid)
	err = dredis.Del(key)
	if err != nil {
		resp.SetResult(dhttp.ERROR_CODE_HTTP_TOKEN_AUTHENTICATION_FAILED, err.Error())
		dlog.Error(dhttp.ERROR_CODE_HTTP_TOKEN_AUTHENTICATION_FAILED, err.Error())
		return resp
	}

	return resp
}
