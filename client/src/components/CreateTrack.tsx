import * as React from 'react'
import { Form, Button } from 'semantic-ui-react'
import Auth from '../auth/Auth'
import { getUploadUrl, uploadFile, createTrack } from '../api/tracks-api'

enum UploadState {
  NoUpload,
  FetchingPresignedUrl,
  UploadingFile,
}

interface CreateTrackProps {
  match: {
    params: {
      trackId: string
    }
  }
  auth: Auth,
  history: any
}

interface CreateTrackState {
  name: string,
  description: string,
  file: any,
  uploadState: UploadState
}

export class CreateTrack extends React.PureComponent<CreateTrackProps, CreateTrackState> {
  state: CreateTrackState = {
    name: "",
    description: "",
    file: undefined,
    uploadState: UploadState.NoUpload
  }

  handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (!files) return

    this.setState({
      file: files[0]
    })
  }

  onTextAreaChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const { name, value } = event.currentTarget
    this.setState(prevState => ({
      ...prevState,
      [name]: value,
    }))
  }

  onInputchange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.currentTarget
    this.setState(prevState => ({
      ...prevState,
      [name]: value,
    }))
  }

  handleSubmit = async (event: React.SyntheticEvent) => {
    event.preventDefault()
    try {
      if (!this.state.file) {
        alert('File should be selected')
        return
      }

      this.setUploadState(UploadState.FetchingPresignedUrl)
      // 1. Creates the track persinting the future URL
      // 2. Returns the uploadUrl in that response
      const newTrack = await createTrack(this.props.auth.getIdToken(), {
        name: this.state.name,
        description: this.state.description
      })

      const uploadUrl = newTrack.signedUploadUrl
      // await getUploadUrl(this.props.auth.getIdToken(), this.props.match.params.trackId)

      this.setUploadState(UploadState.UploadingFile)
      await uploadFile(uploadUrl, this.state.file)

      // TODO: Could present some success alert.
      this.props.history.push(`/`)
    } catch (e) {
      alert('Could not create the track: ' + e.message)
    } finally {
      this.setUploadState(UploadState.NoUpload)
    }
  }

  setUploadState(uploadState: UploadState) {
    this.setState({
      uploadState
    })
  }

  render() {
    return (
      <div>
        <h1>Show us what you got!</h1>

        <Form onSubmit={this.handleSubmit}>
          <Form.Field>
            <label>File</label>
            <input
              type="file"
              accept="audio/*"
              placeholder="Song to upload"
              onChange={this.handleFileChange}
            />
          </Form.Field>


          <Form.Field>
            <label>Name</label>
            <input placeholder='Name' name="name" onChange={this.onInputchange} />
          </Form.Field>

          <Form.Field>
            <label>About</label>
            <textarea placeholder="About" name="description" onChange={this.onTextAreaChange} />
          </Form.Field>


          {this.renderButton()}
        </Form>
      </div>
    )
  }

  renderButton() {

    return (
      <div>
        {this.state.uploadState === UploadState.FetchingPresignedUrl && <p>Uploading image metadata</p>}
        {this.state.uploadState === UploadState.UploadingFile && <p>Uploading file</p>}
        <Button
          loading={this.state.uploadState !== UploadState.NoUpload}
          type="submit"
        >
          Send
        </Button>
      </div>
    )
  }
}
