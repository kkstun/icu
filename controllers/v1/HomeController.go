/**
 * Created by IntelliJ IDEA.
 * User: Kevin
 * Date: 2020/5/22
 * Time: 15:37
 */
package v1

import (
	"app/icu/CUtil"
	"app/icu/config"
	"app/icu/models"
	"app/icu/service"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"path"
	"strings"
)

type HomeController struct {
	BaseController
}

/**
 * 随机头像
 */
func (this *HomeController) RandAvatar() {
	randIndex := CUtil.RandInt(1,config.AvatarNums) - 1

	this.JsonResponse(1,"OK",map[string]int{"avatar":randIndex})
}

/**
 * 修改用户信息
 */
func (this *HomeController) Modify ()  {
	avatar,err := this.GetInt("avatar",0)
	if err != nil {
		this.JsonResponse(-1,err.Error(),nil)
	}

	if avatar < 0 || avatar > config.AvatarNums {
		this.JsonResponse(-1,"系统头像、别瞎搞",nil)
	}

	nick := this.GetString("nick","")
	size := len([]rune(nick))

	if nick == "" || size > 20 {
		this.JsonResponse(-1,"昵称不能为空且不大于20个字符",nil)
	}

	email := this.GetString("email","")

	if email != "" {
		valid  := validation.Validation{}
		result := valid.Email(email,"Email")
		if !result.Ok {
			this.JsonResponse(-1,"邮箱格式不正确",nil)
		}
	}

	sex,err := this.GetInt("sex",0)
	if err != nil {
		this.JsonResponse(-1,err.Error(),nil)
	}

	if sex < 0 || sex > 2 {
		this.JsonResponse(-1,"性别参数错误",nil)
	}


	age,err := this.GetInt("age",0)
	if err != nil {
		this.JsonResponse(-1,err.Error(),nil)
	}

	if age < 1 || age > 80 {
		this.JsonResponse(-1,"年龄范围1~80",nil)
	}

	_,err = models.UpdateUserInfo(this.GetUserId(true),orm.Params{
		"avatar":avatar,
		"nick"  :nick,
		"email" :email,
		"sex"   :sex,
		"age"   :age,
	})

	if err != nil {
		this.JsonResponse(-1,err.Error(),nil)
	}

	this.JsonResponse(1,"OK",nick)
}

//发送聊天文件
func (this *HomeController) Upload() {

	//发送消息
	hasJoin,AClient := service.HasJoinRoom(this.GetUserId(true))
	if !hasJoin {
		this.JsonResponse(-1,"您已不在聊天室",nil)
	}

	key := "img"

	file, handler, err := this.GetFile(key)
	if err != nil {
		this.JsonResponse(-1,err.Error(),nil)
	}
	defer file.Close()

	fileExt := path.Ext(handler.Filename)

	fileExt  = strings.ToLower(fileExt)

	canUpload := false

	for _,val := range config.OnlyUploadImgExt {
		if val == fileExt {
			canUpload = true
			break
		}
	}

	if !canUpload {
		exts := strings.Join(config.OnlyUploadImgExt," , ")
		this.JsonResponse(-1,"仅允许上传"+exts+"格式图片",nil)
	}

	err,link := this.UpLoadFileToServer(key,"chatFiles",config.UploadSize)
	if err != nil {
		this.JsonResponse(-1,err.Error(),nil)
	}

	//发图片
	AClient.SendChatFile(link,models.EventSendChatImg)

	this.JsonResponse(1,"OK",map[string]string{"link":link})
}

//分享网络音频
func (this *HomeController) NetWorkAudio() {

	link := this.GetString("link")

	//发送消息
	hasJoin,AClient := service.HasJoinRoom(this.GetUserId(true))
	if !hasJoin {
		this.JsonResponse(-1,"您已不在聊天室",nil)
	}

	//发图片
	AClient.SendChatFile(link,models.EventSendChatAudio)

	this.JsonResponse(1,"OK",nil)
}

//分享网络视频
func (this *HomeController) NetWorkVideo() {

	link := this.GetString("link")

	//发送消息
	hasJoin,AClient := service.HasJoinRoom(this.GetUserId(true))
	if !hasJoin {
		this.JsonResponse(-1,"您已不在聊天室",nil)
	}

	//发图片
	AClient.SendChatFile(link,models.EventSendChatVedio)

	this.JsonResponse(1,"OK",nil)
}