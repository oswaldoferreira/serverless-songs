import { CustomAuthorizerEvent, CustomAuthorizerResult } from 'aws-lambda'
import 'source-map-support/register'

import { verify, decode } from 'jsonwebtoken'
// import { createLogger } from '../../utils/console'
import { Jwt } from '../Jwt'
import { JwtPayload } from '../JwtPayload'
import * as jwksClient from 'jwks-rsa'

// const console = createLogger('auth')

const jwksUrl = 'https://dev-todo-app.us.auth0.com/.well-known/jwks.json'

export const handler = async (
  event: CustomAuthorizerEvent
): Promise<CustomAuthorizerResult> => {
  console.log(`Authorizing a user ${event.authorizationToken}`)
  try {
    const jwtToken = await verifyToken(event.authorizationToken)
    console.log(`User was authorized ${jwtToken.sub}`)

    return {
      principalId: jwtToken.sub,
      policyDocument: {
        Version: '2012-10-17',
        Statement: [
          {
            Action: 'execute-api:Invoke',
            Effect: 'Allow',
            Resource: '*'
          }
        ]
      },
      context: {
        userId: jwtToken.sub,
      },
    }
  } catch (e) {
    console.log('User not authorized', { error: e.message })

    return {
      principalId: 'user',
      policyDocument: {
        Version: '2012-10-17',
        Statement: [
          {
            Action: 'execute-api:Invoke',
            Effect: 'Deny',
            Resource: '*'
          }
        ]
      }
    }
  }
}

async function verifyToken(authHeader: string): Promise<JwtPayload> {
  const token = getToken(authHeader)
  const jwt: Jwt = decode(token, { complete: true }) as Jwt
  const headerKid = jwt.header.kid

  const client = jwksClient({ jwksUri: jwksUrl })
  const key = await client.getSigningKeyAsync(headerKid)
  const cert = key.getPublicKey();

  return verify(token, cert, { algorithms: ['RS256'] }) as JwtPayload
}

function getToken(authHeader: string): string {
  if (!authHeader) throw new Error('No authentication header')

  if (!authHeader.toLowerCase().startsWith('bearer '))
    throw new Error('Invalid authentication header')

  const split = authHeader.split(' ')
  const token = split[1]

  return token
}