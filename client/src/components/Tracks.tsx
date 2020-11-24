import dateFormat from 'dateformat'
import { History } from 'history'
import update from 'immutability-helper'
import * as React from 'react'
import {
  Button,
  Checkbox,
  Divider,
  Grid,
  Header,
  Icon,
  Input,
  Image,
  Loader
} from 'semantic-ui-react'

import { createTrack, deleteTrack, getTracks, patchTrack } from '../api/tracks-api'
import Auth from '../auth/Auth'
import { Track } from '../types/Track'

interface TracksProps {
  auth: Auth
  history: History
}

interface TracksState {
  tracks: Track[]
  newTrackName: string
  loadingTracks: boolean
}

export class Tracks extends React.PureComponent<TracksProps, TracksState> {
  state: TracksState = {
    tracks: [],
    newTrackName: '',
    loadingTracks: true
  }

  handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    this.setState({ newTrackName: event.target.value })
  }

  onTrackDelete = async (trackId: string) => {
    try {
      await deleteTrack(this.props.auth.getIdToken(), trackId)
      this.setState({
        tracks: this.state.tracks.filter(track => track.trackId != trackId)
      })
    } catch {
      alert('Track deletion failed')
    }
  }

  async componentDidMount() {
    try {
      const tracks = await getTracks(this.props.auth.getIdToken())
      this.setState({
        tracks,
        loadingTracks: false
      })
    } catch (e) {
      alert(`Failed to fetch tracks: ${e.message}`)
    }
  }

  render() {
    return (
      <div>
        {this.renderTracks()}
      </div>
    )
  }

  renderTracks() {
    if (this.state.loadingTracks) {
      return this.renderLoading()
    }

    return this.renderTracksList()
  }

  renderHeader() {
    if (this.state.tracks.length === 0) {
      return ""
    }

    return (
        <Header as="h1">Tracks</Header>
    )
  }

  renderLoading() {
    return (
      <Grid.Row>
        <Loader indeterminate active inline="centered">
          Loading tracks
        </Loader>
      </Grid.Row>
    )
  }

  renderTracksList() {
    return (
      <div>
        {this.renderHeader()}
        <Grid padded>
          {this.state.tracks.map((track, pos) => {
            return (
              <Grid.Row key={track.trackId}>
                <Grid.Column width={6} verticalAlign="middle">
                  <audio controls preload="none">
                    <source src={track.trackUrl} type="audio/mpeg" />
                  </audio> 
                </Grid.Column>

                <Grid.Column width={2} floated="left">
                  {track.name}
                </Grid.Column>

                <Grid.Column width={4} floated="left">
                  {track.description}
                </Grid.Column>

                <Grid.Column width={1} floated="right">
                  <Button
                    icon
                    color="red"
                    onClick={() => this.onTrackDelete(track.trackId)}
                  >
                    <Icon name="delete" />
                  </Button>
                </Grid.Column>
              </Grid.Row>
            )
          })}
        </Grid>

      </div>
    )
  }
}
