/* tslint:disable */
/**
 * This file was automatically generated by json-schema-to-typescript.
 * DO NOT MODIFY IT BY HAND. Instead, modify the source JSONSchema file,
 * and run json-schema-to-typescript to regenerate this file.
 */

declare type ImageV1 = ImageFileV1[] | ImageFileV1;
declare type ImageFileV1 =
  | {
      /**
       * relative to bundle root
       */
      path: string;
      /**
       * (at 1x)
       */
      width?: number;
      /**
       * (at 1x)
       */
      height?: number;
      /**
       * MIME type
       */
      type?: string;
      /**
       * pixel density descriptor, followed by 'x' (e.g. "2x")
       */
      density?: string;
      [k: string]: unknown;
    }
  | string;
declare type VideoV1 = VideoFileV1[] | VideoFileV1;
declare type VideoFileV1 =
  | {
      /**
       * relative to bundle root
       */
      path: string;
      width?: number;
      height?: number;
      /**
       * MIME type
       */
      type?: string;
      [k: string]: unknown;
    }
  | string;

declare interface Artbundle {
  /**
   * Artwork title (to be kept short)
   */
  title: string;
  /**
   * Artwork author (to be kept short)
   */
  author: string;
  /**
   * Artwork license
   */
  license: string;
  /**
   * Artwork creation date
   */
  date?: string;
  /**
   * description of the artwork
   */
  description?: string;
  /**
   * legal notices, copyrights…
   */
  attribution?: string;
  preview?: ImageV1;
  /**
   * List of artwork descriptions. The first supported one will be used.
   */
  artwork: Artwork[];
  [k: string]: unknown;
}
declare interface Artwork {
  /**
   * artwork format
   */
  format: "text" | "image" | "video" | "html";
  /**
   * supported orientation
   */
  orientation: "horizontal" | "vertical" | "both";
  color?: {
    [k: string]: unknown;
  };
  /**
   * configuration of the artwork, corresponding to the specified 'format'
   */
  configuration: ArtworkTextV1 | ArtworkImageV1 | ArtworkHtmlV1 | ArtworkVideoV1;
  [k: string]: unknown;
}
/**
 * 'text' artwork configuration
 */
declare interface ArtworkTextV1 {
  text: string;
  "text-color"?: {
    [k: string]: unknown;
  };
  "background-color"?: {
    [k: string]: unknown;
  };
  [k: string]: unknown;
}
/**
 * 'image' artwork configuration
 */
declare interface ArtworkImageV1 {
  image: ImageV1;
  mode?: "cover" | "contain" | "repeat";
  "background-color"?: {
    [k: string]: unknown;
  };
  [k: string]: unknown;
}
/**
 * 'html' artwork configuration
 */
declare interface ArtworkHtmlV1 {
  /**
   * path to HTML folder to serve
   */
  src?: string;
  [k: string]: unknown;
}
/**
 * 'video' artwork configuration
 */
declare interface ArtworkVideoV1 {
  video: VideoV1;
  mode?: "cover" | "contain" | "repeat";
  "background-color"?: {
    [k: string]: unknown;
  };
  [k: string]: unknown;
}