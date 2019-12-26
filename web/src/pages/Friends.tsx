import * as React from "react"
import { FriendsListURL } from "../constants/api"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
interface friend {
	id: number
}
export const Friends = () => {
	const [friends, setFriends] = React.useState<friend[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	React.useEffect(() => {
		const fetchFriends = async () => {
			try {
				const res = await fetch(FriendsListURL)
				if (!res.ok) {
					const err: Error = await res.json()
					throw new Error(err.message)
				}

				const data: { data: friend[] } = await res.json()
				console.log(data)
				setFriends(data.data)
			} catch (err) {
				console.error(err)
				setErr(err.toString())
			}
		}
		fetchFriends()
	})
	return (
		<div>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			<h1>Friends</h1>
			{!friends && <p>No data</p>}
			{friends &&
				friends.map(friend => {
					return <li key={friend.id}>{friend.id}</li>
				})}
		</div>
	)
}
