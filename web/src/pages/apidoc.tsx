import * as React from "react"
interface Props { }
export const APIDocumentation = (props: Props) => {
    return (<div>

        <h2>Private Routes</h2>
        <ul>
            <li>Name: signOut</li>
            <li>Method: Post</li>
            <li>URL: /api/auth/sign_out</li>
            <li>Comment: destroys cookie</li>
        </ul>

        <ul>
            <li>Name: check</li>
            <li>Method: Get</li>
            <li>URL: /api/auth/check</li>
            <li>Comment: returns signed in user object</li>
        </ul>

        <ul>
            <li>Name: setPassword</li>
            <li>Method: Post</li>
            <li>URL: /api/auth/set_password</li>
            <li>Comment: Not done</li>
        </ul>

        <ul>
            <li>Name: userJWT</li>
            <li>Method: Get</li>
            <li>URL: /api/auth/jwt</li>
            <li>Comment: returns token</li>
        </ul>

        <ul>
            <li>Name: blob</li>
            <li>Method: Get</li>
            <li>URL: /api/blobs/{`{blob_id}`}</li>
            <li>Comment: download images</li>
        </ul>

        <ul>
            <li>Name: userList</li>
            <li>Method: Get</li>
            <li>URL: /api/users/list</li>
            <li>Comment: list users (admin only)</li>
        </ul>

        <ul>
            <li>Name: userImpersonate</li>
            <li>Method: Post</li>
            <li>URL: /api/users/impersonate/{`{user_id}`}</li>
            <li>Comment: impersonate users (admin only)</li>
        </ul>

        <ul>
            <li>Name: integrationsList</li>
            <li>Method: Get</li>
            <li>URL: /api/integrations/list</li>
            <li>Comment: list current user's VRC integrations</li>
        </ul>

        <ul>
            <li>Name: integrationsAddUsername</li>
            <li>Method: Post</li>
            <li>URL: /api/integrations/add_username</li>
            <li>Comment: add VRC integration {`{username, password}`}</li>
        </ul>

        <ul>
            <li>Name: integrationUpdateFriends</li>
            <li>Method: Post</li>
            <li>URL: /api/integrations/{`{integration_id}`}/update_friends</li>
            <li>Comment: update cached friends (done every 5 min auto)</li>
        </ul>

        <ul>
            <li>Name: integrationsDelete</li>
            <li>Method: Post</li>
            <li>URL: /api/integrations/{`{integration_id}`}/delete</li>
            <li>Comment: </li>
        </ul>

        <ul>
            <li>Name: attendanceList</li>
            <li>Method: Get</li>
            <li>URL: /api/integrations/{`{integration_id}`}/attendance/{`{teacher_id}`}/list</li>
            <li>Comment: </li>
        </ul>

        <ul>
            <li>Name: friendList</li>
            <li>Method: Get</li>
            <li>URL: /api/integrations/{`{integration_id}`}/friends/list</li>
            <li>Comment:</li>
        </ul>
        <ul>
            <li>Name: friendRefresh</li>
            <li>Method: Post</li>
            <li>URL: /api/integrations/{`{integration_id}`}/friends/refresh</li>
            <li>Comment:</li>
        </ul>
        <ul>
            <li>Name: friendPromote</li>
            <li>Method: Post</li>
            <li>URL: /api/integrations/{`{integration_id}`}/friends/{`{friend_id}`}/promote</li>
            <li>Comment: Move friend to teacher status</li>
        </ul>

        <ul>
            <li>Name: friendDemote</li>
            <li>Method: Post</li>
            <li>URL: /api/integrations/{`{integration_id}`}/friends/{`{friend_id}`}/demote</li>
            <li>Comment: Move teacher to student status</li>
        </ul>


        <h2>Public Routes</h2>
        <ul>
            <li>Name: signIn</li>
            <li>Method: Post</li>
            <li>URL: /api/auth/sign_in</li>
            <li>Comment: {`{email, password}`}</li>
        </ul>

        <ul>
            <li>Name: signUp</li>
            <li>Method: Post</li>
            <li>URL: /api/auth/sign_up</li>
            <li>Comment: {`{email, password}`}</li>
        </ul>

        <ul>
            <li>Name: metrics</li>
            <li>Method: Get</li>
            <li>URL: /api/metrics</li>
            <li>Comment:</li>
        </ul>
    </div>)
}