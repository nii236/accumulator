import * as React from "react"
import { FriendsListURL, TeachersListURL } from "../constants/api"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
interface teacher {
	id: number
}
export const Teachers = () => {
	const [teachers, setTeachers] = React.useState<teacher[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	React.useEffect(() => {
		const fetchTeachers = async () => {
			try {
				const res = await fetch(TeachersListURL)
				if (!res.ok) {
					const err: Error = await res.json()
					throw new Error(err.message)
				}

				const data: { data: teacher[] } = await res.json()
				console.log(data)
				setTeachers(data.data)
			} catch (err) {
				console.error(err)
				setErr(err.toString())
			}
		}
		fetchTeachers()
	})
	return (
		<div>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			<h1>Teachers</h1>
			{!teachers && <p>No data</p>}
			{teachers &&
				teachers.map(teacher => {
					return <li key={teacher.id}>{teacher.id}</li>
				})}
		</div>
	)
}
