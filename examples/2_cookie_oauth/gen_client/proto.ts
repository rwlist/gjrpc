interface Transport {
    request(method: string, params: unknown): Promise<unknown>;
}

export class Auth {
    constructor(private transport: Transport) {}
    

    async Oauth(): Promise<OAuthResponse> {
        return await this.transport.request("auth.oauth", null) as OAuthResponse;
    }

    async Status(): Promise<AuthStatus> {
        return await this.transport.request("auth.status", null) as AuthStatus;
    }
}

export interface OAuthResponse {
    RedirectURL: string
}

export interface AuthStatus {
    Email: string
}
