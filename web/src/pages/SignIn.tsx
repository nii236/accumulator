import * as React from "react"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { useForm } from "react-hook-form"
import { Input } from "baseui/input"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faUserAlt, faKey } from "@fortawesome/free-solid-svg-icons"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
export const SignIn = () => {
	const { register, setValue, handleSubmit, errors, setError } = useForm<{
		email: string
		password: string
	}>()
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [err, setErr] = React.useState<string | null>(null)
	const signIn = async (data: { email: string; password: string }) => {
		setThinking(true)
		try {
			const res = await fetch("/api/auth/sign_in", { method: "POST", body: JSON.stringify(data) })
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

	const onSubmit = handleSubmit(({ email, password }) => {
		setError([])
		setErr(null)
		signIn({ email, password })
	})
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	return (
		<form onSubmit={onSubmit}>
			{err && <Notification kind={KIND.negative}>{err}</Notification>}
			{errors.email && <Notification kind={KIND.negative}>This field is required</Notification>}
			<Input startEnhancer={<FontAwesomeIcon icon={faUserAlt} />} name="email" placeholder="email" inputRef={register({ required: true })} />
			{errors.password && <Notification kind={KIND.negative}>This field is required</Notification>}
			<Input startEnhancer={<FontAwesomeIcon icon={faKey} />} name="password" type="password" placeholder="password" inputRef={register({ required: true })} />
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
