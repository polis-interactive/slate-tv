
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


void main(void)
{

    float t = time / 2.0;


	vec3 color = vec3(0.0);

    vec2 uv = gl_FragCoord.xy / resolution.xy;

    vec2 uv_grid = floor(vec2(uv.x * 14.0, uv.y * 14.0));


    float separation = (0.5 * 0.5) * 2.0 * 3.14159 + 1.0;

    float mod_offset = mod(uv_grid.x + uv_grid.y, 2.0) * separation;


    float pct = - pow(sin(mod_offset + t / 2.0), 2.0) + t / 10.0;

    pct = sin(pct);

    color = hsb2rgb(vec3(pct,1.0,1.0));

    gl_FragColor = vec4(color, 1.0);
}