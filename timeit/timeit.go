package timeit

// The timeit package provides a function `Timeit` for timing grequests`s HTTP resquest
import (
	"errors"
	"time"

	"github.com/levigross/grequests"
)

type Method string

var (
	POST   Method = "POST"
	GET    Method = "GET"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

func TimeIt(method Method, url string, ro *grequests.RequestOptions) (milliseconds int64, resp *grequests.Response, err error) {
	start := time.Now()
	switch method {
	case POST:
		resp, err = grequests.Post(url, ro)
	case GET:
		resp, err = grequests.Get(url, ro)
	case PUT:
		resp, err = grequests.Put(url, ro)
	case DELETE:
		resp, err = grequests.Delete(url, ro)
	case PATCH:
		resp, err = grequests.Patch(url, ro)
	default:
		return 0, nil, errors.New("method not supported")
	}
	du := time.Since(start)
	return du.Milliseconds(), resp, err
}
