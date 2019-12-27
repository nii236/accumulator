import * as React from "react"
import { Spinner } from "baseui/spinner"
import { Error } from "../types/api"
import { RouteComponentProps } from "react-router-dom"
import { friend } from "./Friends"
import { faMonument } from "@fortawesome/free-solid-svg-icons"
import * as moment from "moment"

interface attendanceRecord {
	location: string
	integration_id: number
	teacher_id: number
	friend_id: number
	timestamp: number
}
interface Props extends RouteComponentProps<{ integration_id: string; teacher_id: string }> {}
export const Attendance = (props: Props) => {
	const [attendances, setAttendance] = React.useState<attendanceRecord[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [friends, setFriends] = React.useState<friend[] | null>(null)
	React.useEffect(() => {
		fetchAttendance()
	}, [])
	React.useEffect(() => {
		fetchFriends()
	}, [])
	const fetchFriends = async () => {
		setThinking(true)
		try {
			const res = await fetch(`/api/integrations/${props.match.params.integration_id}/friends/list`)
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: friend[] } = await res.json()
			setFriends(data.data)
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		setThinking(false)
	}
	const fetchAttendance = async () => {
		setThinking(true)
		try {
			const res = await fetch(`/api/integrations/${props.match.params.integration_id}/attendance/${props.match.params.teacher_id}/list`)
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: attendanceRecord[] } = await res.json()
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
	if (!friends) return <p>No data</p>
	const teacher = friends.find(friend => friend.id.toString() === props.match.params.teacher_id)
	if (!teacher) return <p>Teacher not found</p>
	if (!attendances) return <p>No data</p>
	return (
		<div>
			<h1>{teacher.vrchat_display_name}</h1>
			{groupBy(attendances, att => {
				return att.friend_id
			}).map((arr, i) => {
				const thisFriend = friends.find(el => {
					if (!arr[0]) {
						return false
					}
					return el.id === arr[0].friend_id
				})
				let name = "Unknown name"
				if (thisFriend) {
					name = thisFriend.vrchat_display_name
				}
				return (
					<div key={i}>
						<h2>Student: {name}</h2>
						{arr
							.sort((a, b) => {
								if (a.timestamp === b.timestamp) {
									return 0
								}
								if (a.timestamp > b.timestamp) {
									return -1
								}
								if (a.timestamp < b.timestamp) {
									return 1
								}
								return 0
							})
							.map((el, j) => {
								if (el.teacher_id !== teacher.id) {
									return
								}
								if (el.integration_id.toString() !== props.match.params.integration_id) {
									return
								}
								return (
									<div key={`${i}-${j}`}>
										<ul>
											<li>
												{moment.unix(el.timestamp).toISOString()} - {el.location}
											</li>
										</ul>
									</div>
								)
							})}
					</div>
				)
			})}
		</div>
	)
}
function groupBy<T, K>(list: T[], getKey: (item: T) => K) {
	const map = new Map<K, T[]>()
	list.forEach(item => {
		const key = getKey(item)
		const collection = map.get(key)
		if (!collection) {
			map.set(key, [item])
		} else {
			collection.push(item)
		}
	})
	return Array.from(map.values())
}
