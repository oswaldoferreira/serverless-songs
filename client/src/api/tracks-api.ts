import { apiEndpoint } from '../config'
import { Track } from '../types/Track';
import { CreateTrackRequest } from '../types/CreateTrackRequest';
import { CreateTrackResponse } from '../types/CreateTrackResponse';
import Axios from 'axios'
import { UpdateTrackRequest } from '../types/UpdateTrackRequest';

export async function getTracks(idToken: string): Promise<Track[]> {

  const response = await Axios.get(`${apiEndpoint}/tracks`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${idToken}`
    },
  })
  return response.data
}

export async function createTrack(
  idToken: string,
  newTrack: CreateTrackRequest
): Promise<CreateTrackResponse> {
  const response = await Axios.post(`${apiEndpoint}/tracks`,  JSON.stringify(newTrack), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${idToken}`
    }
  })
  return response.data
}

export async function patchTrack(
  idToken: string,
  trackId: string,
  updatedTrack: UpdateTrackRequest
): Promise<void> {
  await Axios.patch(`${apiEndpoint}/tracks/${trackId}`, JSON.stringify(updatedTrack), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${idToken}`
    }
  })
}

export async function deleteTrack(
  idToken: string,
  trackId: string
): Promise<void> {
  await Axios.delete(`${apiEndpoint}/tracks/${trackId}`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${idToken}`
    }
  })
}

export async function getUploadUrl(
  idToken: string,
  trackId: string
): Promise<string> {
  const response = await Axios.post(`${apiEndpoint}/tracks/${trackId}/attachment`, '', {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${idToken}`
    }
  })
  return response.data.uploadUrl
}

export async function uploadFile(uploadUrl: string, file: Buffer): Promise<void> {
  await Axios.put(uploadUrl, file)
}
