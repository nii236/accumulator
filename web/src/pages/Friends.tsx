import * as React from "react"
import { FriendsListURL } from "../constants/api"
import { Spinner } from "baseui/spinner"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { RouteComponentProps, Redirect } from "react-router-dom"
import { Card, StyledBody, StyledAction } from "baseui/card"
import { Button } from "baseui/button"
import { FlexGrid, FlexGridItem } from "baseui/flex-grid"
import { BlockProps } from "baseui/block"
import { Root } from "baseui/toast"
import { UI } from "../controllers/ui"
export interface friend {
	id: number
	is_teacher: boolean
	vrchat_id: string
	vrchat_username: string
	vrchat_display_name: string
	vrchat_avatar_image_url: string
	vrchat_avatar_thumbnail_image_url: string
	vrchat_location: string
}
interface Props extends RouteComponentProps<{ integration_id: string }> {}
export const Friends = (props: Props) => {
	const [friends, setFriends] = React.useState<friend[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [redirect, setRedirect] = React.useState<string | null>(null)
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
	const promoteToTeacher = async (vrchat_friend_id: string) => {
		try {
			const res = await fetch(`/api/integrations/${props.match.params.integration_id}/friends/${vrchat_friend_id}/promote`, { method: "POST" })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
	}
	const demoteToStudent = async (vrchat_friend_id: string) => {
		try {
			const res = await fetch(`/api/integrations/${props.match.params.integration_id}/friends/${vrchat_friend_id}/demote`, { method: "POST" })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
	}

	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	if (redirect) {
		return <Redirect to={redirect} push />
	}
	if (!friends) return <p>No data</p>
	const itemProps: BlockProps = {
		// backgroundColor: "mono300",
		// height: "scale1000",
		display: "flex",
		// alignItems: "center",
		// justifyContent: "center",
	}
	return (
		<div>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			<h1>Friends</h1>
			<FlexGrid flexWrap={true} flexDirection={"row"} flexGridColumnCount={3} flexGridColumnGap="scale800" flexGridRowGap="scale800">
				{friends &&
					friends.map(friend => {
						return (
							<FlexGridItem key={friend.id} {...itemProps}>
								<FriendItemContainer
									key={friend.id}
									friend={friend}
									demoteToStudent={demoteToStudent}
									promoteToTeacher={promoteToTeacher}
									fetchFriends={fetchFriends}
									setRedirect={setRedirect}
									integration_id={props.match.params.integration_id}
								/>
							</FlexGridItem>
						)
					})}
			</FlexGrid>
		</div>
	)
}
interface FriendItemContainerProps {
	friend: friend
	demoteToStudent: (friendID: string) => void
	promoteToTeacher: (friendID: string) => void
	fetchFriends: () => void
	setRedirect: (url: string) => void
	integration_id: string
}
const FriendItemContainer = (props: FriendItemContainerProps) => {
	const role = props.friend.is_teacher ? "teacher" : "student"
	const location = props.friend.vrchat_location
	const demote = () => {
		props.demoteToStudent(props.friend.vrchat_id)
	}
	const promote = () => {
		props.promoteToTeacher(props.friend.vrchat_id)
	}
	const fetch = () => {
		props.fetchFriends()
	}
	const redirectToAttendance = () => {
		props.setRedirect(`/integrations/${props.integration_id}/attendance/${props.friend.id}`)
	}
	return (
		<FriendCard
			role={role}
			location={location}
			demote={demote}
			promote={promote}
			fetch={fetch}
			redirectToAttendance={redirectToAttendance}
			headerImageURL={props.friend.vrchat_avatar_image_url}
			title={props.friend.vrchat_display_name}
		/>
	)
}

interface ItemProps {
	role: "student" | "teacher"
	location: string
	demote: () => void
	promote: () => void
	fetch: () => void
	redirectToAttendance: () => void
	headerImageURL: string
	title: string
}
const FriendCard = (props: ItemProps) => {
	const ui = UI.useContainer()
	return (
		<Card overrides={{ Root: { style: { width: "100%", height: "100%" } } }} headerImage={props.headerImageURL} title={props.title}>
			<StyledBody>
				<p>
					<em>{props.role}</em>
				</p>
				<p>
					<em>{props.location}</em>
				</p>
			</StyledBody>
			<StyledAction>
				{props.role == "teacher" && (
					<>
						<Button
							onClick={async () => {
								ui.startThinking()
								await props.demote()
								await props.fetch()
								ui.stopThinking()
							}}
							overrides={{
								BaseButton: { style: { width: "100%" } },
							}}>
							Set as student
						</Button>

						<Button
							onClick={() => props.redirectToAttendance()}
							overrides={{
								BaseButton: { style: { width: "100%" } },
							}}>
							View attendances
						</Button>
					</>
				)}
				{props.role !== "teacher" && (
					<Button
						onClick={async () => {
							ui.startThinking()
							await props.promote()
							await props.fetch()
							ui.stopThinking()
						}}
						overrides={{
							BaseButton: { style: { width: "100%" } },
						}}>
						Set as teacher
					</Button>
				)}
			</StyledAction>
		</Card>
	)
}
