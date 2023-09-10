package string_util

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/satori/go.uuid"
)

func UUID() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func ToIntSafe(str string) int {
	v, e := strconv.Atoi(str)
	if nil != e {
		return 0
	}
	return v
}

func ToInt64Safe(str string) int64 {
	v, e := strconv.ParseInt(str, 10, 64)
	if nil != e {
		return 0
	}
	return v
}

func ToUint64Safe(str string) uint64 {
	v, e := strconv.ParseUint(str, 10, 64)
	if nil != e {
		return 0
	}
	return v
}

func Uint64ToStr(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func CommasToMap(base string, sep string) map[string]interface{} {
	ret := make(map[string]interface{})
	if "" != base && "" != sep {
		kvs := strings.Split(base, sep)
		for _, kv := range kvs {
			temp := strings.Split(kv, "=")
			if len(temp) < 2 {
				continue
			}
			if temp[0] == "" {
				continue
			}
			ret[temp[0]] = temp[1]
		}
	}
	return ret
}

func ToJsonBytes(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if nil != err {
		return nil
	}
	return bytes
}

func ToJsonString(v interface{}) string {
	bytes, err := json.Marshal(v)
	if nil != err {
		return ""
	}
	return string(bytes)
}

func ToJsonIndent(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if nil != err {
		return ""
	}
	return string(bytes)
}

func UrlValuesToQueryString(base string, parameters url.Values) string {
	if len(parameters) == 0 {
		return base
	}

	if !strings.Contains(base, "?") {
		base += "?"
	}

	if strings.HasSuffix(base, "?") || strings.HasSuffix(base, "&") {
		base += parameters.Encode()
	} else {
		base += "&" + parameters.Encode()
	}

	return base
}

func MapToQueryString(base string, parameters map[string]interface{}) string {
	if len(parameters) == 0 {
		return base
	}

	exist := false
	if strings.Contains(base, "?") {
		exist = true
	}
	var buffer bytes.Buffer
	buffer.WriteString(base)
	for k, v := range parameters {
		var temp string
		if !exist {
			temp = "?" + k + "=" + ToString(v)
			exist = true
		} else {
			temp = "&" + k + "=" + ToString(v)
		}
		buffer.WriteString(temp)
	}
	return buffer.String()
}

func ToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func IsEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSHA256(plaintext string, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(plaintext))
	return fmt.Sprintf("%X", hash.Sum(nil))
}

func HmacMD5(plaintext string, key string) string {
	hash := hmac.New(md5.New, []byte(key))
	hash.Write([]byte(plaintext))
	return fmt.Sprintf("%X", hash.Sum(nil))
}

func Camel2Case(name string) string {
	buffer := new(bytes.Buffer)
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteString(strconv.FormatInt(int64(r), 10))
		}
	}
	return buffer.String()
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	name = strings.Replace(name, " ", "", -1)
	return Lcfirst(name)
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func ToFloat64Safe(str string) float64 {
	v, e := strconv.ParseFloat(str, 64)
	if nil != e {
		return 0
	}
	return v
}

func ToUint32(str string) (uint32, error) {
	v, e := strconv.ParseUint(str, 10, 32)
	if nil != e {
		return 0, e
	}
	return uint32(v), nil
}

func ToUint32Safe(str string) uint32 {
	v, e := strconv.ParseUint(str, 10, 32)
	if nil != e {
		return 0
	}
	return uint32(v)
}
