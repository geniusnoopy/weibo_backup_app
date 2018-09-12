package dto

type WeiboListQryRespDto struct {
	Ok int
	Data *data
}

type data struct {
	CardlistInfo *cardlistInfo
	Cards *[]card
	ShowAppTips int
	Scheme string
}

type cardlistInfo struct {
	Containerid string
	V_p int
	Show_style int
	Total int
	Page int
}

type card struct {
	Card_type int
	Itemid string
	Scheme string
	Mblog *mblog
	Show_type int
}

type mblog struct {
	Created_at string
	Id string
	Idstr string
	Mid string
	Can_edit bool
	Text string
	Source string
	Favorited bool
	Is_paid bool
	Mblog_vip_type int
	User *user
	Retweeted_status *retweeted_status
	Reposts_count int
	Comments_count int
	Attitudes_count int
	Pending_approval_count int
	IsLongText bool
	Hide_flag int
	Visible *visible
	Mblogtype int
	More_info_type int
	Content_auth int
	Edit_config *edit_config
	Weibo_position int
	Raw_text string
	Bid string
	Pics *[]pic
}

type pic struct {
    Pid string
    Url string
    Size string
    Large *largePic
}

type largePic struct {
    Size string
    Url string
}

type user struct {
	Id int
	Screen_name string
	Profile_image_url string
	Profile_url string
	Statuses_count int
	Verified bool
	Verified_type int
	Close_blue_v bool
	Description string
	Gender string
	Mbtype int
	Urank int
	Mbrank int
	Follow_me bool
	Following bool
	Followers_count int
	Follow_count int
	Cover_image_phone string
	Avatar_hd string
	Like bool
	Like_me bool
	Badge *badge
}

type badge struct {
	Anniversary int
	Unread_pool int
	Unread_pool_ext int
	User_name_certificate int
}

type retweeted_status struct {

}

type visible struct {

}

type edit_config struct {
	Edited bool
}