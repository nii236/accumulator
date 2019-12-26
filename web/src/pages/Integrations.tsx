import * as React from "react"
import { TeachersListURL, IntegrationsAddUsernameURL, IntegrationsAddAPIKeyURL, IntegrationsListURL } from "../constants/api"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { useForm } from "react-hook-form"

const AddVRChatAPIKey = () => {
	const { register, setValue, handleSubmit, errors } = useForm<{
		apiKey: string
		authToken: string
	}>()
	const [err, setErr] = React.useState<string | null>(null)
	const submit = async (apiKey: string, authToken: string) => {
		try {
			const res = await fetch(IntegrationsAddAPIKeyURL, { method: "POST", body: JSON.stringify({ apiKey, authToken }) })
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
	const onSubmit = handleSubmit(({ apiKey, authToken }) => {
		submit(apiKey, authToken)
	})

	return (
		<form onSubmit={onSubmit}>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			{errors.apiKey && <Notification kind={KIND.negative}>This field is required</Notification>}
			<label>API Key</label>
			<input name="apiKey" ref={register({ required: true })} />
			{errors.authToken && <Notification kind={KIND.negative}>This field is required</Notification>}
			<label>Auth Token</label>
			<input name="authToken" ref={register({ required: true })} />
			<input type="submit" />
		</form>
	)
}

const AddVRChatUsername = () => {
	const { register, setValue, handleSubmit, errors } = useForm<{
		username: string
		password: string
	}>()
	const [err, setErr] = React.useState<string | null>(null)
	const submit = async (username: string, password: string) => {
		try {
			const res = await fetch(IntegrationsAddUsernameURL, { method: "POST" })
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
			<input type="submit" />
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
			<AddVRChatAPIKey />
			<AddVRChatUsername />
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			<h1>Integrations</h1>
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
