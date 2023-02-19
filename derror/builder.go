package derror

type (
	errorBuilder struct {
		code    int
		desc    string
		message string
		trace   string
	}

	ErrorBuilder interface {
		Code(code int) ErrorBuilder
		Message(message string) ErrorBuilder
		Desc(desc string) ErrorBuilder
		Build() ServerError
	}
)

var _ ErrorBuilder = (*errorBuilder)(nil)

func (b *errorBuilder) Code(code int) ErrorBuilder {
	b.code = code
	return b
}

func (b *errorBuilder) Message(message string) ErrorBuilder {
	b.message = message
	return b
}

func (b *errorBuilder) Desc(desc string) ErrorBuilder {
	b.desc = desc
	return b
}

func (b *errorBuilder) Build() ServerError {
	return &serverError{
		message: b.message,
		code:    b.code,
		desc:    b.desc,
	}
}
