package constants

// Attachment Types
const (
	AttachmentTypeMergedOhReport = "merged_oh_report"
	AttachmentTypeRx             = "rx"
	AttachmentTypeVitals         = "vitals"
	AttachmentTypeReport         = "report"
	AttachmentTypeVisitDocument  = "visit_document"
	AttachmentTypeTestDocument   = "test_document"
	AttachmentTypeCpDocument     = "cp_document"
)

const (
	FolderPathVisitDocument = "visit_document"
	FolderPathTestDocument  = "test_document"
)

// Attachment Labels
const (
	AttachmentLabelLabDocs       = "lab_docs"
	AttachmentLabelNonLabDocs    = "non_lab_docs"
	AttachmentLabelOutsourceDocs = "outsource_docs"
)

// Attachment Extensions
const (
	AttachmentExtensionPDF  = "pdf"
	AttachmentExtensionPNG  = "png"
	AttachmentExtensionJPG  = "jpg"
	AttachmentExtensionJPEG = "jpeg"
)

var ATTACHMENT_LABELS = []string{
	AttachmentLabelLabDocs,
	AttachmentLabelNonLabDocs,
	AttachmentLabelOutsourceDocs,
}

var ATTACHMENT_TYPES = []string{
	AttachmentTypeMergedOhReport,
	AttachmentTypeRx,
	AttachmentTypeVitals,
	AttachmentTypeReport,
	AttachmentTypeVisitDocument,
	AttachmentTypeTestDocument,
	AttachmentTypeCpDocument,
}

var ALLOWED_ATTACHMENT_EXTESNIONS = []string{
	AttachmentExtensionPDF,
	AttachmentExtensionPNG,
	AttachmentExtensionJPG,
	AttachmentExtensionJPEG,
}

var IMAGE_EXTENSIONS = []string{
	AttachmentExtensionPNG,
	AttachmentExtensionJPG,
	AttachmentExtensionJPEG,
}
