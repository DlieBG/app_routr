import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { ServiceServiceService } from '../service-service.service';
import { DevService } from '../types';

@Component({
  selector: 'app-service-info',
  templateUrl: './service-info.component.html',
  styleUrls: ['./service-info.component.scss']
})
export class ServiceInfoComponent implements OnInit {

  service$!: Observable<DevService>;
  service!: DevService;

  constructor(private serviceService: ServiceServiceService) { }

  ngOnInit(): void {
    this.getService();
  }

  getService() {
    this.service$ = this.serviceService.getService();
    this.service$.subscribe((data) => {
      this.service = data;
    });
  }

}
