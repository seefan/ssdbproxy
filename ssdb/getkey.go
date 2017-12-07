package ssdb

func GetKey(resp []interface{}) string {
	switch resp[0] {
	case "set", "setx", "setnx", "expire", "ttl", "get", "getset", "del", "incr", "exists", "getbit", "setbit", "bitcount", "substr", "strlen", "hset", "hget", "hdel", "hincr", "hexists", "hsize", "hkeys", "hgetall", "hscan", "hrscan", "hclear", "multi_hset", "multi_hget", "multi_hdel", "zset", "zget", "zdel", "zincr", "zsize", "zexists", "zkeys", "zscan", "zrscan", "zrank, zrrank", "zrange, zrrange", "zclear", "zcount", "zsum", "zavg", "zremrangebyrank", "zremrangebyscore", "zpop_front", "zpop_back", "multi_zset", "multi_zget", "multi_zdel", "qsize", "qclear", "qfront", "qback", "qget", "qset", "qrange", "qslice", "qpush", "qpush_front", "qpush_back", "qpop", "qpop_front", "qpop_back", "qtrim_front", "qtrim_back":
		return resp[1].(string)
	}
	return ""
}
