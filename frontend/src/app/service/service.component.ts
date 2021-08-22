import { Component, Input, OnInit } from '@angular/core';
import { ServiceServiceService } from '../service-service.service';
import { DevRoute } from '../types';

@Component({
  selector: 'app-service',
  templateUrl: './service.component.html',
  styleUrls: ['./service.component.scss']
})
export class ServiceComponent implements OnInit {

  @Input() route!: DevRoute;
  @Input() i!: number;

  newEndpoint!: string;

  constructor(private serviceService: ServiceServiceService) { }

  ngOnInit(): void {
  }

  endpointChange() {
    this.serviceService.postRoute(this.i, this.route.Endpoint).subscribe((data) => {
      this.route = data;
    });
  }

  endpointAdd() {
    this.serviceService.postRoute(this.i, this.newEndpoint).subscribe((data) => {
      this.route = data;
      this.newEndpoint = '';
    });
  }

}
