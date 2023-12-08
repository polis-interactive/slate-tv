
#ifdef GL_ES

precision highp float;
#define IN varying
#define OUT out
#define TEXTURE texture2D

#else

#define IN in
#define OUT out
#define TEXTURE texture

#endif

uniform float time;
uniform vec2 resolution;

float rand(float n){return fract(sin(n) * 43758.5453123);}

float noise(float p){
    float fl = floor(p);
    float fc = fract(p);
    return mix(rand(fl), rand(fl + 1.0), fc);
}


vec3 rgb2hsb( in vec3 c ){
    vec4 K = vec4(0.0, -1.0 / 3.0, 2.0 / 3.0, -1.0);
    vec4 p = mix(vec4(c.bg, K.wz),
    vec4(c.gb, K.xy),
    step(c.b, c.g));
    vec4 q = mix(vec4(p.xyw, c.r),
    vec4(c.r, p.yzx),
    step(p.x, c.r));
    float d = q.x - min(q.w, q.y);
    float e = 1.0e-10;
    return vec3(abs(q.z + (q.w - q.y) / (6.0 * d + e)),
    d / (q.x + e),
    q.x);
}


vec3 hsb2rgb( in vec3 c ){
    vec3 rgb = clamp(abs(mod(c.x*6.0+vec3(0.0,4.0,2.0),
    6.0)-3.0)-1.0,
    0.0,
    1.0 );
    rgb = rgb*rgb*(3.0-2.0*rgb);
    return c.z * mix(vec3(1.0), rgb, c.y);
}


void main(){

    vec2 uv = gl_FragCoord.xy / resolution.xy;

    float t = 4.0 * (time + 1.0);

    vec3 color = vec3(0.0);


    vec2 uv_grid = floor(vec2(uv.x * 14.0, uv.y * 14.0)) / vec2(14.0, 14.0);

    float pct = -uv_grid.x * 0.5 * 2.0 + noise(time / 4.0 + 32.48 * rand(uv_grid.y)) * 0.2 - pow(sin(uv_grid.y * 3.0 + time / 2.0), 2.0) + time / 10.0;

    color = hsb2rgb(vec3(pct,1.0,1.0));

    gl_FragColor = vec4(color, 1.0);


}
