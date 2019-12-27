import * as React from "react"

import { HeaderNavigation, ALIGN, StyledNavigationItem as NavigationItem, StyledNavigationList as NavigationList } from "baseui/header-navigation"
import { StyledLink as Link } from "baseui/link"
import { Button } from "baseui/button"
import { Redirect } from "react-router-dom"
import { Spinner } from "baseui/spinner"
import { useUI, UI } from "../controllers/ui"
import { Modal, ModalHeader, ModalBody, ModalFooter, ModalButton } from "baseui/modal"
export const Nav = () => {
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [redirect, setRedirect] = React.useState<string | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const { startThinking } = UI.useContainer()
	const [isOpen, setIsOpen] = React.useState(false)

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
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	return (
		<HeaderNavigation>
			<NavigationList $align={ALIGN.left}>
				<NavigationItem>
					<Link href="/">Home</Link>
				</NavigationItem>
				<NavigationItem>
					<div
						style={{ textDecoration: "underline" }}
						onClick={() => {
							setIsOpen(true)
						}}>
						API Keys
					</div>
					<React.Fragment>
						<Modal onClose={() => setIsOpen(false)} isOpen={isOpen}>
							<ModalHeader>Hello world</ModalHeader>
							<ModalBody>
								Proin ut dui sed metus pharetra hend rerit vel non mi. Nulla ornare faucibus ex, non facilisis nisl. Maecenas aliquet mauris ut tempus.
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
