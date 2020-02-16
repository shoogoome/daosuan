package elasticsearch

// 新建产品搜索配置
func NewProductSearchConfig(from, size int) Config {
	return Config{
		from: from,
		size: size,
		fields: []string{"name", "description", "details", "additional"},
		source: []string{"id", "name", "update_time", "create_time",
			   			 "description", "cover", "status", "star"},
	}
}
