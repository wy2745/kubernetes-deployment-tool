package locust


/*
{"errors": [],
 "stats": [
 {"median_response_time": 6,
  "min_response_time": 6,
  "current_rps": 0.5,
  "name": "/index.html",
  "num_failures": 0,
  "max_response_time": 6,
  "avg_content_length": 612,
  "avg_response_time": 6.0,
  "method": "GET",
  "num_requests": 10},
  {"median_response_time": 6,
  "min_response_time": 6,
  "current_rps": 0.5,
  "name": "Total",
  "num_failures": 0,
  "max_response_time": 6,
  "avg_content_length": 612,
  "avg_response_time": 6.0,
  "method": null,
 "num_requests": 10}
 ],
 "fail_ratio": 0.0,
 "slave_count": 2,
 "state": "stopped", "user_count": 0, "total_rps": 0.5}
 */

type requestStat struct {
	Errors     []interface{} `json:"errors"`
	Stats      []interface{}  `json:"stats"`
	User_count int       `json:"user_count"`
}
