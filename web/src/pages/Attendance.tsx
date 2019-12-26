import * as React from "react"
import { Spinner } from "baseui/spinner"
import { Error } from "../types/api"
import { RouteComponentProps } from "react-router-dom"
interface attendance {
	timestamp: string
	integration_id: string
	friend_id: string
	teacher_id: string
	location: string
}
interface Props extends RouteComponentProps<{ integration_id: string }> {}
export const Attendance = (props: Props) => {
	const [attendances, setAttendance] = React.useState<attendance[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const [thinking, setThinking] = React.useState<boolean>(false)
	React.useEffect(() => {
		fetchAttendance()
	}, [])
	const fetchAttendance = async () => {
		setThinking(true)
		try {
			const res = await fetch(`/api/integrations/${props.match.params.integration_id}/attendance/list`)
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: attendance[] } = await res.json()
			setAttendance(data.data)
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		setThinking(false)
	}
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	if (!attendances) return <p>No data</p>
	return (
		<div>
			{attendances.map(attendance => {
				return (
					<div key={`${attendance.timestamp}${attendance.friend_id}`}>
						<ul>
							<li>{attendance.timestamp}</li>
							<li>{attendance.friend_id}</li>
							<li>{attendance.integration_id}</li>
							<li>{attendance.location}</li>
							<li>{attendance.teacher_id}</li>
						</ul>
					</div>
				)
			})}
		</div>
	)
}
