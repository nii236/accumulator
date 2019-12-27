import * as React from "react"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { useForm } from "react-hook-form"
import { Input } from "baseui/input"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faUserAlt, faKey } from "@fortawesome/free-solid-svg-icons"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { FlexGrid, FlexGridItem } from "baseui/flex-grid"
import { Card } from "baseui/card"
import { BlockProps } from "baseui/block"
import { RouteComponentProps, Redirect } from "react-router-dom"

interface Props extends RouteComponentProps {}
export const SignUp = (props: Props) => {
	const [redirect, setRedirect] = React.useState<string | null>(null)
	const { register, setValue, handleSubmit, errors, setError } = useForm<{
		email: string
		password: string
	}>()
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [err, setErr] = React.useState<string | null>(null)
	const [success, setSuccess] = React.useState<boolean | null>(false)
	const signUp = async (data: { email: string; password: string }) => {
		setThinking(true)
		try {
			const res = await fetch("/api/auth/sign_up", { method: "POST", body: JSON.stringify(data) })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}
			setSuccess(true)
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		setThinking(false)
	}

	const onSubmit = handleSubmit(({ email, password }) => {
		setError([])
		setErr(null)
		signUp({ email, password })
	})
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	const itemProps: BlockProps = {
		alignItems: "center",
		justifyContent: "center",
	}
	if (redirect) {
		return <Redirect to={redirect} push />
	}
	return (
		<Card title="Create account" overrides={{ Root: { style: { width: "500px", marginTop: "50px", marginLeft: "auto", marginRight: "auto" } } }}>
			{success && (
				<Notification overrides={{ Body: { style: { marginLeft: "auto", marginRight: "auto" } } }} kind={KIND.positive}>
					Sign up successful
				</Notification>
			)}
			{!success && (
				<form onSubmit={onSubmit}>
					{err && <Notification kind={KIND.negative}>{err}</Notification>}
					{errors.email && <Notification kind={KIND.negative}>This field is required</Notification>}
					<Input startEnhancer={<FontAwesomeIcon icon={faUserAlt} />} name="email" placeholder="email" inputRef={register({ required: true })} />
					{errors.password && <Notification kind={KIND.negative}>This field is required</Notification>}
					<Input
						startEnhancer={<FontAwesomeIcon icon={faKey} />}
						name="password"
						type="password"
						placeholder="password"
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
			)}
			<hr />
			<Button
				onClick={() => {
					setRedirect("/")
				}}
				overrides={{
					BaseButton: { style: { width: "100%" } },
				}}>
				Back
			</Button>
		</Card>
	)
}
