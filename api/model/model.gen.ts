// Generate model file for angular typescript


export class Wikibook {
  uuid: string;
  subject: string;
  model: string;
  title: string;
  pages: number;
  volumes: Volume[];
}

export class Volume {
  title: string;
  chapters: Chapter[];
}

export class Chapter {
  title: string;
  articles: Page[];
}

export class Page {
  id: number;
  title: string;
}

export class StatusResponse {
  status: string[];
}

export class Void {
}

export class OrderRequest {
  subject: string;
  model: string;
}

export class OrderResponse {
  uuid: string;
}

export class OrderStatusRequest {
  uuid: string;
}

export class OrderStatusResponse {
  status: string;
  wikibook_uuid: string;
}

export class GetWikibookRequest {
  uuid: string;
}

export class GetWikibookResponse {
  wikibook: Wikibook;
}

export class ListWikibookRequest {
}

export class ListWikibookResponse {
  wikibooks: Wikibook[];
}

export class DownloadWikibookRequest {
  uuid: string;
  format: string;
}

