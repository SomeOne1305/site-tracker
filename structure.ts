type Project = {
	id: string
	project_token: string
	name: string
	description: string
	owner: User
	paths: Path[]
}

type Path = {
	id: string
	path: string // Unique
	project: Project
	visit_count: number
	visitors: Visitor[]
}

type Visitor = {
	id: string
	ip_address: string
	user_agent: string
	visit_time: Date
	country: string
}

type User = {
	id: string
	first_name: string
	last_name: string
	email: string
	is_verified: boolean
	verification_code: number
	password: string
	projects: Project[]
	sessions: Session[]
}

type Session = {
	id: string
	refresh_token: string
	expires_at: Date
	created_at: Date
}
