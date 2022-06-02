package controller

import (
	"BytesDanceProject/model"
	"BytesDanceProject/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction implement follow or unfollow the user
func RelationAction(c *gin.Context) {
	//接收post请求携带的信息
	p := new(model.RelationAction)
	if err := c.ShouldBindJSON(p);err != nil{
		//参数错误
		// 类型断言，判断这个错误是否由binding引发的
		_,ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvalidErr)
			return
		}
		//走到这里说明是由binding引发的
		ResponseError(c,CodeNotAccordStandard)
		return
	}
	//业务处理
	service.RelationAction(p)
	ResponseSuccess(c,CodeSuccess,nil)
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	// 接收参数（GET)
	p := &model.FollowListRE{}
	if err := c.ShouldBindQuery(p);err != nil{
		//参数错误
		_,ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvalidErr)
			return
		}
		//走到这里说明是由binding引发的
		ResponseError(c,CodeNotAccordStandard)
		return
	}
	// 逻辑处理
	data ,err := service.GetFollowList(p)
	if err != nil {
		ResponseError(c,CodeServerBusy)
		return
	}
	userdata := make([]User,len(data))
	for idx,tdata := range data{
		userdata[idx].Id = tdata.ID
		userdata[idx].Name = tdata.UserName
		userdata[idx].FollowCount = tdata.FollowCount
		userdata[idx].FollowerCount = tdata.FollowerCount
		userdata[idx].IsFollow = tdata.IsFollow
	}
	ResponseSuccessWithData(c,UserListResponse{
		Response: Response{
			StatusCode: 200,
		},
		UserList: userdata,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	// 接收参数（GET)
	p := &model.FollowListRE{}
	if err := c.ShouldBindQuery(p);err != nil{
		//参数错误
		_,ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvalidErr)
			return
		}
		//走到这里说明是由binding引发的
		ResponseError(c,CodeNotAccordStandard)
		return
	}
	// 逻辑处理
	data ,err := service.GetFollowerList(p)
	if err != nil {
		ResponseError(c,CodeServerBusy)
		return
	}
	userdata := make([]User,len(data))
	for idx,tdata := range data{
		userdata[idx].Id = tdata.ID
		userdata[idx].Name = tdata.UserName
		userdata[idx].FollowCount = tdata.FollowCount
		userdata[idx].FollowerCount = tdata.FollowerCount
		userdata[idx].IsFollow = tdata.IsFollow
	}
	ResponseSuccessWithData(c,UserListResponse{
		Response: Response{
			StatusCode: 200,
		},
		UserList: userdata,
	})
}