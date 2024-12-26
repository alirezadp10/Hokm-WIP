package errors

type ValidationError struct {
    StatusCode int
    Message    string
    Details    interface{}
}
