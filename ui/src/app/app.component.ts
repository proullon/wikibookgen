import { Component } from '@angular/core';
import { BreakpointObserver } from '@angular/cdk/layout';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass']
})
export class AppComponent {
  title = 'ui';
  isMobile = false;

  constructor(mobileBreakPointObserver: BreakpointObserver) {
   const layoutChange =  mobileBreakPointObserver.observe('(max-width: 768px)');

   layoutChange.subscribe((result) => {
     this.isMobile = result.matches;
     console.log(this.isMobile);
   });
  }
}
