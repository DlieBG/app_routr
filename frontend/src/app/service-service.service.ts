import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { DevRoute, DevService } from './types';

@Injectable({
  providedIn: 'root'
})
export class ServiceServiceService {

  constructor(private httpClient: HttpClient) { }

  public getService(): Observable<DevService> {
    return this.httpClient.get<DevService>(environment.apiUrl);
  }

  public postRoute(routeIndex: number, endpoint: string): Observable<DevRoute> {
    return this.httpClient.post<DevRoute>(environment.apiUrl, { route: routeIndex, endpoint });
  }
}
