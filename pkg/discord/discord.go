package discord

type Context struct {
}

func CreateContext() (Context, error) {
	context := Context{}
	return context, nil
}
