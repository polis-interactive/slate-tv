
#ifdef GL_ES

precision highp float;
#define IN attribute
#define OUT varying
#define TEXTURE texture2D

#else

#define IN attribute
#define OUT out
#define TEXTURE texture

#endif

attribute vec4 vPosition;

void main() {
  gl_Position = vPosition;
}
