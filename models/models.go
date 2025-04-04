package models

type EmailRequest struct {
	Sender         string `json:"sender" binding:"required"`
	Title          string `json:"title" binding:"required"`
	ContentBody    string `json:"contentBody" binding:"required"`
	Attachment     string `json:"attachment,omitempty"`
	NameAttachment string `json:"nameAttachment,omitempty"`
}

type EmailResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
