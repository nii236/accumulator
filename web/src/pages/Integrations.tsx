import * as React from "react"
import { TeachersListURL, IntegrationsAddUsernameURL, IntegrationsAddAPIKeyURL, IntegrationsListURL } from "../constants/api"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { useForm } from "react-hook-form"

const AddVRChatUsername = () => {
	const { register, setValue, handleSubmit, errors, setError } = useForm<{
		username: string
		password: string
	}>()
	const [err, setErr] = React.useState<string | null>(null)
	const submit = async (username: string, password: string) => {
		try {
			const res = await fetch(IntegrationsAddUsernameURL, { method: "POST", body: JSON.stringify({ username, password }) })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}

			const data: { data: integration[] } = await res.json()
			console.log(data)
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
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
			<label>Username</label>
			<input name="username" ref={register({ required: true })} />
			{errors.password && <Notification kind={KIND.negative}>This field is required</Notification>}
			<label>Password</label>
			<input name="password" ref={register({ required: true })} />
			<button type="submit">Add</button>
		</form>
	)
}

interface integration {
	id: number
	api_key: string
	auth_token: string
}
export const Integrations = () => {
	const [integrations, setIntegrations] = React.useState<integration[] | null>(null)
	const [err, setErr] = React.useState<string | null>(null)

	React.useEffect(() => {
		const fetchIntegrations = async () => {
			try {
				const res = await fetch(IntegrationsListURL)
				if (!res.ok) {
					const err: Error = await res.json()
					throw new Error(err.message)
				}

				const data: { data: integration[] } = await res.json()
				console.log(data)
				setIntegrations(data.data)
			} catch (err) {
				console.error(err)
				setErr(err.toString())
			}
		}
		fetchIntegrations()
	}, [])
	return (
		<div>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			<h1>Integrations</h1>
			<AddVRChatUsername />
			{!integrations && <p>No data</p>}
			{integrations &&
				integrations.map(integration => {
					return (
						<li key={integration.id}>
							<button>X</button>
							{integration.api_key}:{integration.auth_token}
						</li>
					)
				})}
		</div>
	)
}
