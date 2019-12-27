import * as React from "react"

import { Unstable_StyledTable as Table, Unstable_StyledHeadCell as HeadCell, Unstable_StyledBodyCell as BodyCell } from "baseui/table-grid"
import { Spinner } from "baseui/spinner"
import { Button } from "baseui/button"

interface user {
    id: number
    email: string
    role: string
}

export const Users = () => {
    const [users, setUsers] = React.useState<user[] | null>(null)
    const [err, setErr] = React.useState<string | null>(null)
    const [thinking, setThinking] = React.useState<boolean>(false)
    React.useEffect(() => {
        fetchUsers()
    }, [])
    const impersonateUser = async (user_id: number) => {
        try {
            const res = await fetch(`/api/users/impersonate/${user_id}`, { method: "POST" })
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
    const fetchUsers = async () => {
        setThinking(true)
        try {
            const res = await fetch("/api/users/list")
            if (!res.ok) {
                const err: Error = await res.json()
                throw new Error(err.message)
            }
            const data: { data: user[] } = await res.json()
            console.log(data)
            setUsers(data.data)
        } catch (err) {
            console.error(err)
            setErr(err.toString())
        }
        setThinking(false)
    }
    if (thinking) {
        return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
    }
    return (
        <Table $gridTemplateColumns="4fr 1fr">
            <HeadCell>User</HeadCell>
            <HeadCell>Actions</HeadCell>
            {users && users.map(u => {
                return (
                    <React.Fragment key={`${u.email}`}>
                        <BodyCell>{u.email}</BodyCell>
                        <BodyCell><Button onClick={() => { impersonateUser(u.id) }}>Impersonate</Button></BodyCell>
                    </React.Fragment>
                )
            })}
        </Table>
    )
}