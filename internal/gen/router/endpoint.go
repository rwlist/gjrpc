package router

type endpoint struct {
	path       []string
	methodImpl *methodImpl
}
