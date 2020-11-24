export interface CreateTrackResponse {
    trackId: string
    name: string
    description: string
    createdAt: string
    signedUploadUrl: string
    attachmentUrl?: string
  }