export interface ParseRequest {
  query: string
  dialect?: 'prometheus'
  options?: {
    format?: boolean
    analyze?: boolean
  }
}

export interface ParseResponse {
  formatted?: string
  ast?: ASTNode
  dependencies?: Dependencies
  issues?: Issue[]
  error?: string
}

export interface ASTNode {
  type: string
  span?: [number, number]
  text?: string
  meta?: Record<string, any>
  [key: string]: any
}

export interface Dependencies {
  metrics: string[]
  labels: string[]
  matchers: Matcher[]
  functions: string[]
  aggregations: string[]
}

export interface Matcher {
  label: string
  op: string
  value: string
}

export interface Issue {
  level: 'error' | 'warn' | 'info'
  code: string
  message: string
  span?: [number, number]
}

export interface QueryRangeRequest {
  query: string
  start: number
  end: number
  step: number
  datasource?: string
}

export interface QueryRangeResponse {
  status: string
  data: {
    resultType: string
    result: Series[]
  }
}

export interface Series {
  metric: Record<string, string>
  values: [number, string][]
}
