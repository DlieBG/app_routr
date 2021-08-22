export interface DevService {
    Name: string,
    Description: string,
    Host: string,
    Port: number,
    DevPort: number,
    Routes: DevRoute[]
}

export interface DevRoute {
    Name: string,
    Description: string,
    Path: string,
    Endpoint: string,
    RecentEndpoints: string[]
}