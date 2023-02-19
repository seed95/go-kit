package derror

type (
	builder struct {
		code    int
		desc    string
		message string
		trace   string
	}

	Builder interface {
		Code(code int) Builder
		Message(message string) Builder
		Desc(desc string) Builder
		Build() ServerError
	}
)

var _ Builder = (*builder)(nil)

func NewBuilder() *builder {
	return &builder{}
}

func (b *builder) Code(code int) Builder {
	b.code = code
	return b
}

func (b *builder) Message(message string) Builder {
	b.message = message
	return b
}

func (b *builder) Desc(desc string) Builder {
	b.desc = desc
	return b
}

func (b *builder) Build() ServerError {
	return &serverError{
		message: b.message,
		code:    b.code,
		desc:    b.desc,
	}
}
