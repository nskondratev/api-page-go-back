package handler

type handlerCreateTestCase struct {
	inputData                 string
	responseCode              int
	responseBodyShouldContain string
}

type handlerGetTestCase struct {
	id                        string
	responseCode              int
	responseBodyShouldContain string
}

type handlerDeleteTestCase struct {
	id                        string
	responseCode              int
	responseBodyShouldContain string
}

type handlerUpdateTestCase struct {
	id                        string
	inputData                 string
	responseCode              int
	responseBodyShouldContain string
}

var (
	emptyStr            = ""
	emptyQueryParamsMap = map[string]string{}
)
