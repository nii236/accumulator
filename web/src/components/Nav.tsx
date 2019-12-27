import * as React from "react"

import { HeaderNavigation, ALIGN, StyledNavigationItem as NavigationItem, StyledNavigationList as NavigationList } from "baseui/header-navigation"
import { StyledLink as Link } from "baseui/link"
import { Button } from "baseui/button"
import { Redirect } from "react-router-dom"
import { Spinner } from "baseui/spinner"
import { useUI, UI } from "../controllers/ui"
import { Modal, ModalHeader, ModalBody, ModalFooter, ModalButton } from "baseui/modal"
import { Paragraph1 } from "baseui/typography"
import { useForm } from "react-hook-form"
import { Input } from "baseui/input"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faKey } from "@fortawesome/free-solid-svg-icons"
import { Notification, KIND } from "baseui/notification"
interface Props {
	email: string
	role: string
}
export const ChangePasswordForm = () => {
	const [err, setErr] = React.useState<string | null>(null)
	const [thinking, setThinking] = React.useState(false)
	const { register, setValue, handleSubmit, errors, setError } = useForm<{
		email: string
		password: string
	}>()
	const setPassword = async (data: { password: string }) => {
		try {
			const res = await fetch("/api/auth/set_password", { method: "POST", body: JSON.stringify(data) })
			if (!res.ok) {
				const err: Error = await res.json()
				setErr(err.message)
				throw new Error(err.message)
			}
			return true
		} catch (err) {
			setErr(err.toString())
		}
		return false
	}
	const onSubmit = handleSubmit(async ({ email, password }) => {
		setError([])
		setErr(null)
		try {
			const success = await setPassword({ password })
			if (success) {
				window.location.reload()
			}
		} catch (error) {
			setErr(error)
			console.error(error)
		}
	})
	return (
		<form onSubmit={onSubmit}>
			{err && (
				<Notification overrides={{ Body: { style: { marginLeft: "auto", marginRight: "auto" } } }} kind={KIND.negative}>
					{err}
				</Notification>
			)}
			<Input startEnhancer={<FontAwesomeIcon icon={faKey} />} name="password" type="password" placeholder="password" inputRef={register({ required: true })} />
			<Button
				type="submit"
				overrides={{
					BaseButton: { style: { width: "100%", marginTop: "10px" } },
				}}>
				Save
			</Button>
		</form>
	)
}
export const Nav = (props: Props) => {
	const [modalThinking, setModalThinking] = React.useState<boolean>(false)
	const [redirect, setRedirect] = React.useState<string | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const { startThinking } = UI.useContainer()
	const [isOpen, setIsOpen] = React.useState(false)
	const [setPasswordOpen, setSetPasswordOpen] = React.useState(false)
	const [jwt, setJwt] = React.useState<string | null>(null)

	if (redirect) {
		return <Redirect to={redirect} push />
	}
	const signOut = async () => {
		try {
			const res = await fetch("/api/auth/sign_out", { method: "POST" })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		window.location.href = "/"
	}

	return (
		<HeaderNavigation>
			<NavigationList $align={ALIGN.left}>
				<NavigationItem>
					<Link href="/">Home</Link>
				</NavigationItem>
				<NavigationItem>
					<Link href="/documentation">Docs</Link>
				</NavigationItem>
				<NavigationItem>
					<div
						style={{ textDecoration: "underline" }}
						onClick={async () => {
							setSetPasswordOpen(true)
						}}>
						Change Password
					</div>
					<React.Fragment>
						<Modal onClose={() => setSetPasswordOpen(false)} isOpen={setPasswordOpen}>
							<ModalHeader>Set a new password</ModalHeader>
							<ModalBody>
								<ChangePasswordForm />
							</ModalBody>
						</Modal>
					</React.Fragment>
				</NavigationItem>
				{props.role == "admin" && (
					<NavigationItem>
						<Link href="/users">Users</Link>
					</NavigationItem>
				)}
				<NavigationItem>
					<div
						style={{ textDecoration: "underline" }}
						onClick={async () => {
							setIsOpen(true)
							setModalThinking(true)
							const res = await fetch("/api/auth/jwt")
							const data: { data: string } = await res.json()
							setJwt(data.data)
							setModalThinking(false)
						}}>
						API Keys
					</div>
					<React.Fragment>
						<Modal onClose={() => setIsOpen(false)} isOpen={isOpen}>
							<ModalHeader>Your API token</ModalHeader>
							<ModalBody>
								<Paragraph1>Use the following token in your Authorization Header.</Paragraph1>
								<pre
									style={{
										whiteSpace: "pre-wrap",
										wordWrap: "break-word",
									}}>
									{`Authorization: Bearer ${jwt}`}
								</pre>
							</ModalBody>
							<ModalFooter>
								<ModalButton onClick={() => setIsOpen(false)}>Okay</ModalButton>
							</ModalFooter>
						</Modal>
					</React.Fragment>
				</NavigationItem>
			</NavigationList>
			<NavigationList $align={ALIGN.center} />
			<NavigationList $align={ALIGN.right}>
				<NavigationItem>{props.email}</NavigationItem>
				<NavigationItem>
					<Button
						onClick={async () => {
							startThinking()
							await signOut()
						}}>
						Sign out
					</Button>
				</NavigationItem>
			</NavigationList>
		</HeaderNavigation>
	)
}
