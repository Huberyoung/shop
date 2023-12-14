package handler

import (
	"context"
	"crypto/md5"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop_srvs/global"
	mDatabase "shop_srvs/model/database"
	mUser "shop_srvs/model/user"
	pUser "shop_srvs/proto/user"
)

const (
	Iterations = 10
	SaltLen    = 10
	KeyLen     = 40
)

type UserServer struct {
	pUser.UnsafeUserServer
}

func ModelToResponse(user mUser.User) pUser.UserInfoResponse {
	g := pUser.UserInfoResponse_Gender(user.Gender)
	userInfoRsp := pUser.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		NikeName: user.NickName,
		Gender:   &g,
		Role:     uint64(user.Role),
	}

	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// GetUserList 获取角色列表
func (u *UserServer) GetUserList(ctx context.Context, in *pUser.PageInfo) (*pUser.UserListResponse, error) {
	var users []mUser.User
	result := global.DBEngine.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := pUser.UserListResponse{}
	rsp.Total = result.RowsAffected
	page, pageSize := int(in.PageNum), int(in.PageSize)
	global.DBEngine.Scopes(mDatabase.Paginate(page, pageSize)).Find(&users)
	for _, user := range users {
		us := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &us)
	}
	return &rsp, nil
}

// GetUserByMobile 通过手机号获取角色信息
func (u *UserServer) GetUserByMobile(ctx context.Context, in *pUser.MobileRequest) (*pUser.UserInfoResponse, error) {
	var user mUser.User
	result := global.DBEngine.Where(&mUser.User{Mobile: in.Mobile}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	pu := ModelToResponse(user)
	return &pu, nil

}

// GetUserById 通过id获取角色信息
func (u *UserServer) GetUserById(ctx context.Context, in *pUser.IdRequest) (*pUser.UserInfoResponse, error) {
	var user mUser.User
	result := global.DBEngine.First(&user, in.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	pu := ModelToResponse(user)
	return &pu, nil
}

// CreateUser 创建角色
func (u *UserServer) CreateUser(ctx context.Context, in *pUser.CreateUserRequest) (*pUser.UserInfoResponse, error) {
	var user mUser.User
	//result := global.DBEngine.Where(mUser.User{Mobile: in.Mobile}).First(&user)
	result := global.DBEngine.Session(&gorm.Session{NewDB: true}).Where(&mUser.User{Mobile: in.Mobile}).First(&user)
	if result.RowsAffected > 0 {
		return &pUser.UserInfoResponse{}, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	//  单单执行上面的查询没啥问题，把上面的注释去掉，执行下面的插入也可以成功，
	// 但上面查询，查询不到的话，即走到下面的代码的时候就是报错 报 record not found
	options := password.Options{SaltLen: SaltLen, Iterations: Iterations, KeyLen: KeyLen, HashFunction: md5.New}
	salt, pwd := password.Encode(in.Password, &options)
	b := time.Unix(int64(in.BirthDay), 0).Local()
	user1 := mUser.User{
		Mobile:   in.Mobile,
		Password: fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, pwd),
		NickName: in.Nickname,
		Birthday: &b,
		Gender:   int(in.Gender.Number()),
	}

	result1 := global.DBEngine.Create(&user1)
	if result1.Error != nil {
		return nil, status.Errorf(codes.Internal, result1.Error.Error())
	}

	pu := ModelToResponse(user1)
	return &pu, nil
}

// UpdateUser 更新角色信息
func (u *UserServer) UpdateUser(ctx context.Context, in *pUser.UpdateUserRequest) (*emptypb.Empty, error) {
	var user mUser.User
	result := global.DBEngine.Find(&user, in.Id)
	if result.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "用户不存在")
	}

	b := time.Unix(int64(in.BirthDay), 0)
	user.NickName = in.Nickname
	user.Gender = int(in.Gender.Number())
	user.Birthday = &b

	result = global.DBEngine.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

// CheckPassword 验证密码
func (u *UserServer) CheckPassword(ctx context.Context, in *pUser.PasswordCheckRequest) (*pUser.PasswordCheckResponse, error) {
	split := strings.Split(in.EncryptedPassword, "$")
	options := password.Options{SaltLen: SaltLen, Iterations: Iterations, KeyLen: KeyLen, HashFunction: md5.New}
	verify := password.Verify(in.Password, split[1], split[2], &options)
	return &pUser.PasswordCheckResponse{Success: verify}, nil
}
