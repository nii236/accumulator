import * as React from "react"
import { Modal, ModalHeader, ModalBody, ModalFooter, ModalButton } from "baseui/modal"
import { IntegrationsAddUsernameURL, IntegrationsListURL } from "../constants/api"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { useForm } from "react-hook-form"
import { Redirect, RouteComponentProps } from "react-router-dom"
import { H1, H2, Paragraph1 } from "baseui/typography"
import { Button } from "baseui/button"
import { Card, StyledBody, StyledAction } from "baseui/card"
import { FlexGrid, FlexGridItem } from "baseui/flex-grid"
import { Spinner } from "baseui/spinner"
import { BlockProps } from "baseui/block"
import { Input } from "baseui/input"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faUserAlt, faKey } from "@fortawesome/free-solid-svg-icons"
interface AddProps {
	fetchIntegrations: () => void
	setThinking: (thinking: boolean) => void
}
const AddVRChatUsername = (props: AddProps) => {
	const { register, setValue, handleSubmit, errors, setError } = useForm<{
		username: string
		password: string
	}>()
	const [err, setErr] = React.useState<string | null>(null)

	const submit = async (username: string, password: string) => {
		props.setThinking(true)
		try {
			const res = await fetch(IntegrationsAddUsernameURL, { method: "POST", body: JSON.stringify({ username, password }) })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: integration[] } = await res.json()
			props.fetchIntegrations()
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		props.setThinking(false)
	}
	const onSubmit = handleSubmit(({ username, password }) => {
		setError([])
		setErr(null)
		submit(username, password)
	})

	return (
		<form onSubmit={onSubmit}>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			{errors.username && <Notification kind={KIND.negative}>This field is required</Notification>}
			<Input startEnhancer={<FontAwesomeIcon icon={faUserAlt} />} name="username" placeholder="VRChat username" inputRef={register({ required: true })} />
			{errors.password && <Notification kind={KIND.negative}>This field is required</Notification>}
			<Input
				startEnhancer={<FontAwesomeIcon icon={faKey} />}
				name="password"
				type="password"
				placeholder="VRChat password"
				inputRef={register({ required: true })}
			/>
			<Button
				type="submit"
				overrides={{
					BaseButton: { style: { width: "100%" } },
				}}>
				Create
			</Button>
		</form>
	)
}
interface Props extends RouteComponentProps {}
interface integration {
	id: number
	username: string
	api_key: string
	auth_token: string
}
export const Integrations = (props: Props) => {
	const [integrations, setIntegrations] = React.useState<integration[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const [redirect, setRedirect] = React.useState<string | null>(null)
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [removeModalOpen, setRemoveModalOpen] = React.useState(false)
	React.useEffect(() => {
		fetchIntegrations()
	}, [])
	const updateFriends = async (integration_id: number) => {
		setThinking(true)
		try {
			const res = await fetch(`/api/integrations/${integration_id}/update_friends`, { method: "POST" })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: integration[] } = await res.json()
			fetchIntegrations()
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		setThinking(false)
	}
	const fetchIntegrations = async () => {
		setThinking(true)
		try {
			const res = await fetch(IntegrationsListURL)
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: integration[] } = await res.json()
			setIntegrations(data.data)
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		setThinking(false)
	}
	const deleteIntegration = async (integrationID: number) => {
		setThinking(true)
		try {
			const res = await fetch(`/api/integrations/${integrationID}/delete`, { method: "POST" })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		setThinking(false)
	}
	if (redirect) {
		return <Redirect to={redirect} push />
	}
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	const itemProps: BlockProps = {
		// display: "flex",
		// alignItems: "start",
		// justifyItems: "start",
		// alignContent: "start",
		// justifyContent: "start",
	}
	return (
		<div>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			<H1>Integrations</H1>

			<FlexGrid flexDirection={"row"} flexGridColumnCount={3} flexGridColumnGap="scale100" flexGridRowGap="scale100">
				<FlexGridItem {...itemProps}>
					<Card title={"Create new"}>
						<StyledBody>
							<AddVRChatUsername
								fetchIntegrations={fetchIntegrations}
								setThinking={(thinking: boolean) => {
									setThinking(thinking)
								}}
							/>
						</StyledBody>
					</Card>
				</FlexGridItem>
				{integrations &&
					integrations.map(integration => {
						return (
							<FlexGridItem key={integration.id} {...itemProps}>
								<Card key={integration.id} title={integration.username}>
									<StyledBody>
										<small>API Key: {integration.api_key}</small>
									</StyledBody>
									<StyledAction>
										<Button
											onClick={async () => {
												setRemoveModalOpen(true)
											}}
											overrides={{
												BaseButton: { style: { width: "100%" } },
											}}>
											Remove
										</Button>
										<React.Fragment>
											<Modal onClose={() => setRemoveModalOpen(false)} isOpen={removeModalOpen}>
												<ModalHeader>Destroy integration</ModalHeader>
												<ModalBody>
													<Paragraph1>Warning! This is irreversible! You will lose the current tracked attendance for this integration!</Paragraph1>
												</ModalBody>
												<ModalFooter>
													<ModalButton
														onClick={async () => {
															await deleteIntegration(integration.id)
															fetchIntegrations()
															setRemoveModalOpen(false)
														}}>
														Okay
													</ModalButton>
												</ModalFooter>
											</Modal>
										</React.Fragment>
									</StyledAction>
									<StyledAction>
										<Button
											onClick={() => setRedirect(`/integrations/${integration.id}/friends`)}
											overrides={{
												BaseButton: { style: { width: "100%" } },
											}}>
											Friends
										</Button>
									</StyledAction>
									<StyledAction>
										<Button
											onClick={() => updateFriends(integration.id)}
											overrides={{
												BaseButton: { style: { width: "100%" } },
											}}>
											Update Friends
										</Button>
									</StyledAction>
								</Card>
							</FlexGridItem>
						)
					})}
			</FlexGrid>
		</div>
	)
}
