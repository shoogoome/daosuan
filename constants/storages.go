package constants

const (

	RasDataRoot = "../data"

	// 账户
	RasAccount = RasDataRoot + "/account"

	// 机构
	RasInstitution = RasDataRoot + "/institution"

	// 课程
	RasCourse = RasDataRoot + "/course"

	// 资讯模块
	RasInformation = RasDataRoot + "/information"

	// ----- 子模块 -----

	// 账户头像
	RasAccountAvator = RasAccount + "/avator"

	// 机构logo
	RasInstitutionLogo = RasInstitution + "/logo"

	// 机构经营许可
	RasInstitutionLicense = RasInstitution + "/license"

	// 机构附件
	RasInstitutionAttachments = RasInstitution + "/attachments"

	// 课程封面
	RasCourseCover = RasCourse + "/cover"

	// 课程附件
	RasCourseAttachments = RasCourse + "/attachments"

	// 资讯封面
	RasInformationCover = RasInformation + "/cover"

	// 资讯附件
	RasInformationAttachments = RasInformation + "/attachments"


	// nginx静态资源映射
	NginxResourcePath = "/resource_internal"
)

var StorageMapping = map[string]string {
	"account_avator": RasAccountAvator,
	"institution_logo": RasInstitutionLogo,
	"institution_license": RasInstitutionLicense,
	"institution_attachments": RasInstitutionAttachments,
	"course_cover": RasCourseCover,
	"course_attachments": RasCourseAttachments,
	"information_cover": RasInformationCover,
	"information_attachments": RasInformationAttachments,
}

var MimeToExtMapping = map[string]string {
	"jpg": "image/jpeg",
	"jpeg": "image/jpeg",
	"bmp": "image/bmp",
	"png": "image/png",
	"gif": "image/gif",
	"svg": "image/svg",
}
