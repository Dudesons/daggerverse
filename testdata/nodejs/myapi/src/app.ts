import fastify from 'fastify'
import { usersData } from './mockData.js'

const app = fastify({
  logger: process.env.NODE_ENV === 'development',
})

app.get('/users', async () => {
  return usersData
})

export default app
