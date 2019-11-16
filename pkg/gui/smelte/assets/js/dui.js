
(function(l, r) { if (l.getElementById('livereloadscript')) return; r = l.createElement('script'); r.async = 1; r.src = '//' + (window.location.host || 'localhost').split(':')[0] + ':35729/livereload.js?snipver=1'; r.id = 'livereloadscript'; l.head.appendChild(r) })(window.document);
var dui = (function (smelte) {
    'use strict';

    function noop() { }
    const identity = x => x;
    function add_location(element, file, line, column, char) {
        element.__svelte_meta = {
            loc: { file, line, column, char }
        };
    }
    function run(fn) {
        return fn();
    }
    function blank_object() {
        return Object.create(null);
    }
    function run_all(fns) {
        fns.forEach(run);
    }
    function is_function(thing) {
        return typeof thing === 'function';
    }
    function safe_not_equal(a, b) {
        return a != a ? b == b : a !== b || ((a && typeof a === 'object') || typeof a === 'function');
    }
    function validate_store(store, name) {
        if (!store || typeof store.subscribe !== 'function') {
            throw new Error(`'${name}' is not a store with a 'subscribe' method`);
        }
    }
    function subscribe(store, callback) {
        const unsub = store.subscribe(callback);
        return unsub.unsubscribe ? () => unsub.unsubscribe() : unsub;
    }
    function component_subscribe(component, store, callback) {
        component.$$.on_destroy.push(subscribe(store, callback));
    }

    const is_client = typeof window !== 'undefined';
    let now = is_client
        ? () => window.performance.now()
        : () => Date.now();
    let raf = is_client ? cb => requestAnimationFrame(cb) : noop;

    const tasks = new Set();
    let running = false;
    function run_tasks() {
        tasks.forEach(task => {
            if (!task[0](now())) {
                tasks.delete(task);
                task[1]();
            }
        });
        running = tasks.size > 0;
        if (running)
            raf(run_tasks);
    }
    function loop(fn) {
        let task;
        if (!running) {
            running = true;
            raf(run_tasks);
        }
        return {
            promise: new Promise(fulfil => {
                tasks.add(task = [fn, fulfil]);
            }),
            abort() {
                tasks.delete(task);
            }
        };
    }

    function append(target, node) {
        target.appendChild(node);
    }
    function insert(target, node, anchor) {
        target.insertBefore(node, anchor || null);
    }
    function detach(node) {
        node.parentNode.removeChild(node);
    }
    function destroy_each(iterations, detaching) {
        for (let i = 0; i < iterations.length; i += 1) {
            if (iterations[i])
                iterations[i].d(detaching);
        }
    }
    function element(name) {
        return document.createElement(name);
    }
    function svg_element(name) {
        return document.createElementNS('http://www.w3.org/2000/svg', name);
    }
    function text(data) {
        return document.createTextNode(data);
    }
    function space() {
        return text(' ');
    }
    function empty() {
        return text('');
    }
    function listen(node, event, handler, options) {
        node.addEventListener(event, handler, options);
        return () => node.removeEventListener(event, handler, options);
    }
    function attr(node, attribute, value) {
        if (value == null)
            node.removeAttribute(attribute);
        else if (node.getAttribute(attribute) !== value)
            node.setAttribute(attribute, value);
    }
    function children(element) {
        return Array.from(element.childNodes);
    }
    function set_style(node, key, value, important) {
        node.style.setProperty(key, value, important ? 'important' : '');
    }
    function custom_event(type, detail) {
        const e = document.createEvent('CustomEvent');
        e.initCustomEvent(type, false, false, detail);
        return e;
    }

    let stylesheet;
    let active = 0;
    let current_rules = {};
    // https://github.com/darkskyapp/string-hash/blob/master/index.js
    function hash(str) {
        let hash = 5381;
        let i = str.length;
        while (i--)
            hash = ((hash << 5) - hash) ^ str.charCodeAt(i);
        return hash >>> 0;
    }
    function create_rule(node, a, b, duration, delay, ease, fn, uid = 0) {
        const step = 16.666 / duration;
        let keyframes = '{\n';
        for (let p = 0; p <= 1; p += step) {
            const t = a + (b - a) * ease(p);
            keyframes += p * 100 + `%{${fn(t, 1 - t)}}\n`;
        }
        const rule = keyframes + `100% {${fn(b, 1 - b)}}\n}`;
        const name = `__svelte_${hash(rule)}_${uid}`;
        if (!current_rules[name]) {
            if (!stylesheet) {
                const style = element('style');
                document.head.appendChild(style);
                stylesheet = style.sheet;
            }
            current_rules[name] = true;
            stylesheet.insertRule(`@keyframes ${name} ${rule}`, stylesheet.cssRules.length);
        }
        const animation = node.style.animation || '';
        node.style.animation = `${animation ? `${animation}, ` : ``}${name} ${duration}ms linear ${delay}ms 1 both`;
        active += 1;
        return name;
    }
    function delete_rule(node, name) {
        node.style.animation = (node.style.animation || '')
            .split(', ')
            .filter(name
            ? anim => anim.indexOf(name) < 0 // remove specific animation
            : anim => anim.indexOf('__svelte') === -1 // remove all Svelte animations
        )
            .join(', ');
        if (name && !--active)
            clear_rules();
    }
    function clear_rules() {
        raf(() => {
            if (active)
                return;
            let i = stylesheet.cssRules.length;
            while (i--)
                stylesheet.deleteRule(i);
            current_rules = {};
        });
    }

    let current_component;
    function set_current_component(component) {
        current_component = component;
    }

    const dirty_components = [];
    const binding_callbacks = [];
    const render_callbacks = [];
    const flush_callbacks = [];
    const resolved_promise = Promise.resolve();
    let update_scheduled = false;
    function schedule_update() {
        if (!update_scheduled) {
            update_scheduled = true;
            resolved_promise.then(flush);
        }
    }
    function add_render_callback(fn) {
        render_callbacks.push(fn);
    }
    function flush() {
        const seen_callbacks = new Set();
        do {
            // first, call beforeUpdate functions
            // and update components
            while (dirty_components.length) {
                const component = dirty_components.shift();
                set_current_component(component);
                update(component.$$);
            }
            while (binding_callbacks.length)
                binding_callbacks.pop()();
            // then, once components are updated, call
            // afterUpdate functions. This may cause
            // subsequent updates...
            for (let i = 0; i < render_callbacks.length; i += 1) {
                const callback = render_callbacks[i];
                if (!seen_callbacks.has(callback)) {
                    callback();
                    // ...so guard against infinite loops
                    seen_callbacks.add(callback);
                }
            }
            render_callbacks.length = 0;
        } while (dirty_components.length);
        while (flush_callbacks.length) {
            flush_callbacks.pop()();
        }
        update_scheduled = false;
    }
    function update($$) {
        if ($$.fragment !== null) {
            $$.update($$.dirty);
            run_all($$.before_update);
            $$.fragment && $$.fragment.p($$.dirty, $$.ctx);
            $$.dirty = null;
            $$.after_update.forEach(add_render_callback);
        }
    }

    let promise;
    function wait() {
        if (!promise) {
            promise = Promise.resolve();
            promise.then(() => {
                promise = null;
            });
        }
        return promise;
    }
    function dispatch(node, direction, kind) {
        node.dispatchEvent(custom_event(`${direction ? 'intro' : 'outro'}${kind}`));
    }
    const outroing = new Set();
    let outros;
    function group_outros() {
        outros = {
            r: 0,
            c: [],
            p: outros // parent group
        };
    }
    function check_outros() {
        if (!outros.r) {
            run_all(outros.c);
        }
        outros = outros.p;
    }
    function transition_in(block, local) {
        if (block && block.i) {
            outroing.delete(block);
            block.i(local);
        }
    }
    function transition_out(block, local, detach, callback) {
        if (block && block.o) {
            if (outroing.has(block))
                return;
            outroing.add(block);
            outros.c.push(() => {
                outroing.delete(block);
                if (callback) {
                    if (detach)
                        block.d(1);
                    callback();
                }
            });
            block.o(local);
        }
    }
    const null_transition = { duration: 0 };
    function create_in_transition(node, fn, params) {
        let config = fn(node, params);
        let running = false;
        let animation_name;
        let task;
        let uid = 0;
        function cleanup() {
            if (animation_name)
                delete_rule(node, animation_name);
        }
        function go() {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            if (css)
                animation_name = create_rule(node, 0, 1, duration, delay, easing, css, uid++);
            tick(0, 1);
            const start_time = now() + delay;
            const end_time = start_time + duration;
            if (task)
                task.abort();
            running = true;
            add_render_callback(() => dispatch(node, true, 'start'));
            task = loop(now => {
                if (running) {
                    if (now >= end_time) {
                        tick(1, 0);
                        dispatch(node, true, 'end');
                        cleanup();
                        return running = false;
                    }
                    if (now >= start_time) {
                        const t = easing((now - start_time) / duration);
                        tick(t, 1 - t);
                    }
                }
                return running;
            });
        }
        let started = false;
        return {
            start() {
                if (started)
                    return;
                delete_rule(node);
                if (is_function(config)) {
                    config = config();
                    wait().then(go);
                }
                else {
                    go();
                }
            },
            invalidate() {
                started = false;
            },
            end() {
                if (running) {
                    cleanup();
                    running = false;
                }
            }
        };
    }
    function create_out_transition(node, fn, params) {
        let config = fn(node, params);
        let running = true;
        let animation_name;
        const group = outros;
        group.r += 1;
        function go() {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            if (css)
                animation_name = create_rule(node, 1, 0, duration, delay, easing, css);
            const start_time = now() + delay;
            const end_time = start_time + duration;
            add_render_callback(() => dispatch(node, false, 'start'));
            loop(now => {
                if (running) {
                    if (now >= end_time) {
                        tick(0, 1);
                        dispatch(node, false, 'end');
                        if (!--group.r) {
                            // this will result in `end()` being called,
                            // so we don't need to clean up here
                            run_all(group.c);
                        }
                        return false;
                    }
                    if (now >= start_time) {
                        const t = easing((now - start_time) / duration);
                        tick(1 - t, t);
                    }
                }
                return running;
            });
        }
        if (is_function(config)) {
            wait().then(() => {
                // @ts-ignore
                config = config();
                go();
            });
        }
        else {
            go();
        }
        return {
            end(reset) {
                if (reset && config.tick) {
                    config.tick(1, 0);
                }
                if (running) {
                    if (animation_name)
                        delete_rule(node, animation_name);
                    running = false;
                }
            }
        };
    }
    function create_bidirectional_transition(node, fn, params, intro) {
        let config = fn(node, params);
        let t = intro ? 0 : 1;
        let running_program = null;
        let pending_program = null;
        let animation_name = null;
        function clear_animation() {
            if (animation_name)
                delete_rule(node, animation_name);
        }
        function init(program, duration) {
            const d = program.b - t;
            duration *= Math.abs(d);
            return {
                a: t,
                b: program.b,
                d,
                duration,
                start: program.start,
                end: program.start + duration,
                group: program.group
            };
        }
        function go(b) {
            const { delay = 0, duration = 300, easing = identity, tick = noop, css } = config || null_transition;
            const program = {
                start: now() + delay,
                b
            };
            if (!b) {
                // @ts-ignore todo: improve typings
                program.group = outros;
                outros.r += 1;
            }
            if (running_program) {
                pending_program = program;
            }
            else {
                // if this is an intro, and there's a delay, we need to do
                // an initial tick and/or apply CSS animation immediately
                if (css) {
                    clear_animation();
                    animation_name = create_rule(node, t, b, duration, delay, easing, css);
                }
                if (b)
                    tick(0, 1);
                running_program = init(program, duration);
                add_render_callback(() => dispatch(node, b, 'start'));
                loop(now => {
                    if (pending_program && now > pending_program.start) {
                        running_program = init(pending_program, duration);
                        pending_program = null;
                        dispatch(node, running_program.b, 'start');
                        if (css) {
                            clear_animation();
                            animation_name = create_rule(node, t, running_program.b, running_program.duration, 0, easing, config.css);
                        }
                    }
                    if (running_program) {
                        if (now >= running_program.end) {
                            tick(t = running_program.b, 1 - t);
                            dispatch(node, running_program.b, 'end');
                            if (!pending_program) {
                                // we're done
                                if (running_program.b) {
                                    // intro — we can tidy up immediately
                                    clear_animation();
                                }
                                else {
                                    // outro — needs to be coordinated
                                    if (!--running_program.group.r)
                                        run_all(running_program.group.c);
                                }
                            }
                            running_program = null;
                        }
                        else if (now >= running_program.start) {
                            const p = now - running_program.start;
                            t = running_program.a + running_program.d * easing(p / running_program.duration);
                            tick(t, 1 - t);
                        }
                    }
                    return !!(running_program || pending_program);
                });
            }
        }
        return {
            run(b) {
                if (is_function(config)) {
                    wait().then(() => {
                        // @ts-ignore
                        config = config();
                        go(b);
                    });
                }
                else {
                    go(b);
                }
            },
            end() {
                clear_animation();
                running_program = pending_program = null;
            }
        };
    }
    function create_component(block) {
        block && block.c();
    }
    function mount_component(component, target, anchor) {
        const { fragment, on_mount, on_destroy, after_update } = component.$$;
        fragment && fragment.m(target, anchor);
        // onMount happens before the initial afterUpdate
        add_render_callback(() => {
            const new_on_destroy = on_mount.map(run).filter(is_function);
            if (on_destroy) {
                on_destroy.push(...new_on_destroy);
            }
            else {
                // Edge case - component was destroyed immediately,
                // most likely as a result of a binding initialising
                run_all(new_on_destroy);
            }
            component.$$.on_mount = [];
        });
        after_update.forEach(add_render_callback);
    }
    function destroy_component(component, detaching) {
        const $$ = component.$$;
        if ($$.fragment !== null) {
            run_all($$.on_destroy);
            $$.fragment && $$.fragment.d(detaching);
            // TODO null out other refs, including component.$$ (but need to
            // preserve final state?)
            $$.on_destroy = $$.fragment = null;
            $$.ctx = {};
        }
    }
    function make_dirty(component, key) {
        if (!component.$$.dirty) {
            dirty_components.push(component);
            schedule_update();
            component.$$.dirty = blank_object();
        }
        component.$$.dirty[key] = true;
    }
    function init(component, options, instance, create_fragment, not_equal, props) {
        const parent_component = current_component;
        set_current_component(component);
        const prop_values = options.props || {};
        const $$ = component.$$ = {
            fragment: null,
            ctx: null,
            // state
            props,
            update: noop,
            not_equal,
            bound: blank_object(),
            // lifecycle
            on_mount: [],
            on_destroy: [],
            before_update: [],
            after_update: [],
            context: new Map(parent_component ? parent_component.$$.context : []),
            // everything else
            callbacks: blank_object(),
            dirty: null
        };
        let ready = false;
        $$.ctx = instance
            ? instance(component, prop_values, (key, ret, value = ret) => {
                if ($$.ctx && not_equal($$.ctx[key], $$.ctx[key] = value)) {
                    if ($$.bound[key])
                        $$.bound[key](value);
                    if (ready)
                        make_dirty(component, key);
                }
                return ret;
            })
            : prop_values;
        $$.update();
        ready = true;
        run_all($$.before_update);
        // `false` as a special case of no DOM component
        $$.fragment = create_fragment ? create_fragment($$.ctx) : false;
        if (options.target) {
            if (options.hydrate) {
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.l(children(options.target));
            }
            else {
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.c();
            }
            if (options.intro)
                transition_in(component.$$.fragment);
            mount_component(component, options.target, options.anchor);
            flush();
        }
        set_current_component(parent_component);
    }
    class SvelteComponent {
        $destroy() {
            destroy_component(this, 1);
            this.$destroy = noop;
        }
        $on(type, callback) {
            const callbacks = (this.$$.callbacks[type] || (this.$$.callbacks[type] = []));
            callbacks.push(callback);
            return () => {
                const index = callbacks.indexOf(callback);
                if (index !== -1)
                    callbacks.splice(index, 1);
            };
        }
        $set() {
            // overridden by instance, if it has props
        }
    }

    function dispatch_dev(type, detail) {
        document.dispatchEvent(custom_event(type, detail));
    }
    function append_dev(target, node) {
        dispatch_dev("SvelteDOMInsert", { target, node });
        append(target, node);
    }
    function insert_dev(target, node, anchor) {
        dispatch_dev("SvelteDOMInsert", { target, node, anchor });
        insert(target, node, anchor);
    }
    function detach_dev(node) {
        dispatch_dev("SvelteDOMRemove", { node });
        detach(node);
    }
    function listen_dev(node, event, handler, options, has_prevent_default, has_stop_propagation) {
        const modifiers = options === true ? ["capture"] : options ? Array.from(Object.keys(options)) : [];
        if (has_prevent_default)
            modifiers.push('preventDefault');
        if (has_stop_propagation)
            modifiers.push('stopPropagation');
        dispatch_dev("SvelteDOMAddEventListener", { node, event, handler, modifiers });
        const dispose = listen(node, event, handler, options);
        return () => {
            dispatch_dev("SvelteDOMRemoveEventListener", { node, event, handler, modifiers });
            dispose();
        };
    }
    function attr_dev(node, attribute, value) {
        attr(node, attribute, value);
        if (value == null)
            dispatch_dev("SvelteDOMRemoveAttribute", { node, attribute });
        else
            dispatch_dev("SvelteDOMSetAttribute", { node, attribute, value });
    }
    function set_data_dev(text, data) {
        data = '' + data;
        if (text.data === data)
            return;
        dispatch_dev("SvelteDOMSetData", { node: text, data });
        text.data = data;
    }
    class SvelteComponentDev extends SvelteComponent {
        constructor(options) {
            if (!options || (!options.target && !options.$$inline)) {
                throw new Error(`'target' is a required option`);
            }
            super();
        }
        $destroy() {
            super.$destroy();
            this.$destroy = () => {
                console.warn(`Component was already destroyed`); // eslint-disable-line no-console
            };
        }
    }

    function cubicInOut(t) {
        return t < 0.5 ? 4.0 * t * t * t : 0.5 * Math.pow(2.0 * t - 2.0, 3.0) + 1.0;
    }
    function cubicOut(t) {
        const f = t - 1.0;
        return f * f * f + 1.0;
    }

    function fade(node, { delay = 0, duration = 400, easing = identity }) {
        const o = +getComputedStyle(node).opacity;
        return {
            delay,
            duration,
            easing,
            css: t => `opacity: ${t * o}`
        };
    }
    function fly(node, { delay = 0, duration = 400, easing = cubicOut, x = 0, y = 0, opacity = 0 }) {
        const style = getComputedStyle(node);
        const target_opacity = +style.opacity;
        const transform = style.transform === 'none' ? '' : style.transform;
        const od = target_opacity * (1 - opacity);
        return {
            delay,
            duration,
            easing,
            css: (t, u) => `
			transform: ${transform} translate(${(1 - t) * x}px, ${(1 - t) * y}px);
			opacity: ${target_opacity - (od * u)}`
        };
    }
    function draw(node, { delay = 0, speed, duration, easing = cubicInOut }) {
        const len = node.getTotalLength();
        if (duration === undefined) {
            if (speed === undefined) {
                duration = 800;
            }
            else {
                duration = len / speed;
            }
        }
        else if (typeof duration === 'function') {
            duration = duration(len);
        }
        return {
            delay,
            duration,
            easing,
            css: (t, u) => `stroke-dasharray: ${t * len} ${u * len}`
        };
    }

    const inner = `M77.08,2.55c3.87,1.03 6.96,2.58 10.32,4.64c5.93,3.87 10.58,8.51 14.19,14.71c3.87,6.19 5.42,13.16 5.42,20.64c0,7.22 -1.81,14.18 -5.41,20.37c-3.61,6.45 -8.25,11.35 -14.19,14.96c-3.35,2.06 -6.96,3.87 -10.32,4.9c-3.87,1.03 -7.74,1.55 -11.61,1.55v-14.45c6.96,-0.26 13.42,-2.58 19.09,-8c5.67,-5.42 8.51,-11.87 8.51,-19.61c0,-7.74 -2.58,-14.19 -7.99,-19.6c-5.42,-5.42 -11.86,-8 -19.6,-8c-7.74,0 -14.44,2.58 -19.6,8c-5.42,5.42 -8,11.87 -8,19.6l0,85.9c-3.1,-3.1 -7.99,-7.74 -13.93,-13.67v-72.23c0,-3.87 0.52,-7.73 1.55,-11.35c1.03,-3.87 2.58,-7.22 4.64,-10.32c3.87,-5.93 8.52,-10.58 14.71,-14.45c6.19,-3.61 13.16,-5.16 20.64,-5.16c3.87,0 8,0.52 11.61,1.55zM78.37,42.28c0,7.22 -5.93,13.16 -13.15,13.16c-7.48,0.26 -13.16,-5.68 -13.16,-13.16c0,-7.22 5.94,-13.16 13.16,-13.16c7.22,0 13.15,5.93 13.15,13.16zM13.63,37.12l0,69.39c-6.19,-6.19 -11.09,-10.83 -13.93,-13.93l0,-55.46z`;

    /* src/boot/logo/BootLogo.svelte generated by Svelte v3.14.1 */
    const file = "src/boot/logo/BootLogo.svelte";

    function get_each_context(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.char = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    function get_each_context_1(ctx, list, i) {
    	const child_ctx = Object.create(ctx);
    	child_ctx.char = list[i];
    	child_ctx.i = i;
    	return child_ctx;
    }

    // (27:2) {#each 'PLAN 9 FROM CRYPTO SPACE' as char, i}
    function create_each_block_1(ctx) {
    	let span;
    	let t;
    	let span_intro;

    	const block = {
    		c: function create() {
    			span = element("span");
    			t = text(ctx.char);
    			attr_dev(span, "class", "svelte-56x957");
    			add_location(span, file, 27, 3, 570);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, span, anchor);
    			append_dev(span, t);
    		},
    		p: noop,
    		i: function intro(local) {
    			if (!span_intro) {
    				add_render_callback(() => {
    					span_intro = create_in_transition(span, fade, { delay: 33 + ctx.i * 150, duration: 3333 });
    					span_intro.start();
    				});
    			}
    		},
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(span);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block_1.name,
    		type: "each",
    		source: "(27:2) {#each 'PLAN 9 FROM CRYPTO SPACE' as char, i}",
    		ctx
    	});

    	return block;
    }

    // (44:2) {#each 'ParallelCoin' as char, i}
    function create_each_block(ctx) {
    	let span;
    	let t;
    	let span_intro;

    	const block = {
    		c: function create() {
    			span = element("span");
    			t = text(ctx.char);
    			attr_dev(span, "class", "svelte-56x957");
    			add_location(span, file, 44, 3, 1008);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, span, anchor);
    			append_dev(span, t);
    		},
    		p: noop,
    		i: function intro(local) {
    			if (!span_intro) {
    				add_render_callback(() => {
    					span_intro = create_in_transition(span, fade, { delay: 3333 + ctx.i * 150, duration: 999 });
    					span_intro.start();
    				});
    			}
    		},
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(span);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_each_block.name,
    		type: "each",
    		source: "(44:2) {#each 'ParallelCoin' as char, i}",
    		ctx
    	});

    	return block;
    }

    function create_fragment(ctx) {
    	let div0;
    	let div0_outro;
    	let t0;
    	let svg;
    	let g;
    	let path;
    	let path_intro;
    	let g_outro;
    	let t1;
    	let div1;
    	let div1_outro;
    	let t2;
    	let div2;
    	let caption;
    	let t3;
    	let t4;
    	let current;
    	let each_value_1 = "PLAN 9 FROM CRYPTO SPACE";
    	let each_blocks_1 = [];

    	for (let i = 0; i < each_value_1.length; i += 1) {
    		each_blocks_1[i] = create_each_block_1(get_each_context_1(ctx, each_value_1, i));
    	}

    	let each_value = "ParallelCoin";
    	let each_blocks = [];

    	for (let i = 0; i < each_value.length; i += 1) {
    		each_blocks[i] = create_each_block(get_each_context(ctx, each_value, i));
    	}

    	const block = {
    		c: function create() {
    			div0 = element("div");

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].c();
    			}

    			t0 = space();
    			svg = svg_element("svg");
    			g = svg_element("g");
    			path = svg_element("path");
    			t1 = space();
    			div1 = element("div");

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].c();
    			}

    			t2 = space();
    			div2 = element("div");
    			caption = element("caption");
    			t3 = text(ctx.progress);
    			t4 = text("%");
    			attr_dev(div0, "class", "centered marginTop plan svelte-56x957");
    			add_location(div0, file, 25, 0, 445);
    			set_style(path, "stroke", "#cfcfcf");
    			set_style(path, "stroke-width", "1.5");
    			attr_dev(path, "d", inner);
    			attr_dev(path, "class", "svelte-56x957");
    			add_location(path, file, 35, 3, 783);
    			attr_dev(g, "opacity", "0.2");
    			add_location(g, file, 34, 2, 735);
    			attr_dev(svg, "id", "bootlogo");
    			attr_dev(svg, "class", "marginTopBig svelte-56x957");
    			attr_dev(svg, "viewBox", "0 0 108 128");
    			add_location(svg, file, 33, 1, 669);
    			attr_dev(div1, "class", "centered name svelte-56x957");
    			add_location(div1, file, 42, 1, 904);
    			add_location(caption, file, 53, 1, 1164);
    			attr_dev(div2, "class", "progress justifyCenter textCenter txDark svelte-56x957");
    			add_location(div2, file, 52, 0, 1108);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div0, anchor);

    			for (let i = 0; i < each_blocks_1.length; i += 1) {
    				each_blocks_1[i].m(div0, null);
    			}

    			insert_dev(target, t0, anchor);
    			insert_dev(target, svg, anchor);
    			append_dev(svg, g);
    			append_dev(g, path);
    			insert_dev(target, t1, anchor);
    			insert_dev(target, div1, anchor);

    			for (let i = 0; i < each_blocks.length; i += 1) {
    				each_blocks[i].m(div1, null);
    			}

    			insert_dev(target, t2, anchor);
    			insert_dev(target, div2, anchor);
    			append_dev(div2, caption);
    			append_dev(caption, t3);
    			append_dev(caption, t4);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (!current || changed.progress) set_data_dev(t3, ctx.progress);
    		},
    		i: function intro(local) {
    			if (current) return;

    			for (let i = 0; i < each_value_1.length; i += 1) {
    				transition_in(each_blocks_1[i]);
    			}

    			if (div0_outro) div0_outro.end(1);

    			if (!path_intro) {
    				add_render_callback(() => {
    					path_intro = create_in_transition(path, draw, { duration: 3333 });
    					path_intro.start();
    				});
    			}

    			if (g_outro) g_outro.end(1);

    			for (let i = 0; i < each_value.length; i += 1) {
    				transition_in(each_blocks[i]);
    			}

    			if (div1_outro) div1_outro.end(1);
    			current = true;
    		},
    		o: function outro(local) {
    			div0_outro = create_out_transition(div0, fly, { y: -40, duration: 999 });
    			g_outro = create_out_transition(g, fade, { duration: 999 });
    			div1_outro = create_out_transition(div1, fly, { y: -20, duration: 9999 });
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div0);
    			destroy_each(each_blocks_1, detaching);
    			if (detaching && div0_outro) div0_outro.end();
    			if (detaching) detach_dev(t0);
    			if (detaching) detach_dev(svg);
    			if (detaching && g_outro) g_outro.end();
    			if (detaching) detach_dev(t1);
    			if (detaching) detach_dev(div1);
    			destroy_each(each_blocks, detaching);
    			if (detaching && div1_outro) div1_outro.end();
    			if (detaching) detach_dev(t2);
    			if (detaching) detach_dev(div2);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance($$self, $$props, $$invalidate) {
    	let progress = 0;

    	function next() {
    		setTimeout(
    			() => {
    				if (progress === 100) {
    					$$invalidate("progress", progress = 0);
    				}

    				$$invalidate("progress", progress += 1);
    				next();
    			},
    			100
    		);
    	}

    	next();

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		if ("progress" in $$props) $$invalidate("progress", progress = $$props.progress);
    	};

    	return { progress };
    }

    class BootLogo extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance, create_fragment, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "BootLogo",
    			options,
    			id: create_fragment.name
    		});
    	}
    }

    const subscriber_queue = [];
    /**
     * Creates a `Readable` store that allows reading by subscription.
     * @param value initial value
     * @param {StartStopNotifier}start start and stop notifications for subscriptions
     */
    function readable(value, start) {
        return {
            subscribe: writable(value, start).subscribe,
        };
    }
    /**
     * Create a `Writable` store that allows both updating and reading by subscription.
     * @param {*=}value initial value
     * @param {StartStopNotifier=}start start and stop notifications for subscriptions
     */
    function writable(value, start = noop) {
        let stop;
        const subscribers = [];
        function set(new_value) {
            if (safe_not_equal(value, new_value)) {
                value = new_value;
                if (stop) { // store is ready
                    const run_queue = !subscriber_queue.length;
                    for (let i = 0; i < subscribers.length; i += 1) {
                        const s = subscribers[i];
                        s[1]();
                        subscriber_queue.push(s, value);
                    }
                    if (run_queue) {
                        for (let i = 0; i < subscriber_queue.length; i += 2) {
                            subscriber_queue[i][0](subscriber_queue[i + 1]);
                        }
                        subscriber_queue.length = 0;
                    }
                }
            }
        }
        function update(fn) {
            set(fn(value));
        }
        function subscribe(run, invalidate = noop) {
            const subscriber = [run, invalidate];
            subscribers.push(subscriber);
            if (subscribers.length === 1) {
                stop = start(set) || noop;
            }
            run(value);
            return () => {
                const index = subscribers.indexOf(subscriber);
                if (index !== -1) {
                    subscribers.splice(index, 1);
                }
                if (subscribers.length === 0) {
                    stop();
                    stop = null;
                }
            };
        }
        return { set, update, subscribe };
    }

    const bios = readable([], function start(set) {
        const interval = setInterval(() => {
            fetch(`http://127.0.0.1:3999/bios`)
                .then(resp => resp.json())
                .then(data => {
                    set(data);
                });
        }, 100);
        return function stop() {
            clearInterval(interval);
        };
    });

    /* src/boot/Boot.svelte generated by Svelte v3.14.1 */
    const file$1 = "src/boot/Boot.svelte";

    // (9:0) {#if bios.logo}
    function create_if_block(ctx) {
    	let current;
    	const bootlogo = new BootLogo({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(bootlogo.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(bootlogo, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(bootlogo.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(bootlogo.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(bootlogo, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block.name,
    		type: "if",
    		source: "(9:0) {#if bios.logo}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$1(ctx) {
    	let div;
    	let h1;
    	let t1;
    	let current;
    	let if_block = bios.logo && create_if_block(ctx);

    	const block = {
    		c: function create() {
    			div = element("div");
    			h1 = element("h1");
    			h1.textContent = "Boot";
    			t1 = space();
    			if (if_block) if_block.c();
    			attr_dev(h1, "id", "biosMessages");
    			add_location(h1, file$1, 6, 0, 192);
    			attr_dev(div, "class", "fullScreen flx flc justifyEvenly itemsCenter bgDark boot");
    			add_location(div, file$1, 5, 0, 121);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, h1);
    			append_dev(div, t1);
    			if (if_block) if_block.m(div, null);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			if (if_block) if_block.d();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$1.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Boot extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$1, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Boot",
    			options,
    			id: create_fragment$1.name
    		});
    	}
    }

    /* src/components/panels/PanelBalance.svelte generated by Svelte v3.14.1 */

    const file$2 = "src/components/panels/PanelBalance.svelte";

    function create_fragment$2(ctx) {
    	let div7;
    	let div6;
    	let div4;
    	let div2;
    	let div0;
    	let t1;
    	let div1;
    	let span0;
    	let t2;
    	let t3;
    	let div3;
    	let t4;
    	let div5;
    	let small0;
    	let span1;
    	let strong0;
    	let span2;
    	let t6;
    	let small1;
    	let span3;
    	let strong1;
    	let span4;

    	const block = {
    		c: function create() {
    			div7 = element("div");
    			div6 = element("div");
    			div4 = element("div");
    			div2 = element("div");
    			div0 = element("div");
    			div0.textContent = "Balance:";
    			t1 = space();
    			div1 = element("div");
    			span0 = element("span");
    			t2 = text(" DUO");
    			t3 = space();
    			div3 = element("div");
    			t4 = space();
    			div5 = element("div");
    			small0 = element("small");
    			span1 = element("span");
    			span1.textContent = "Pending: ";
    			strong0 = element("strong");
    			span2 = element("span");
    			t6 = space();
    			small1 = element("small");
    			span3 = element("span");
    			span3.textContent = "Transactions: ";
    			strong1 = element("strong");
    			span4 = element("span");
    			attr_dev(div0, "class", "e-card-header-title");
    			add_location(div0, file$2, 32, 4, 589);
    			attr_dev(span0, "v-html", "this.duOSys.status.balance.balance");
    			attr_dev(span0, "class", "svelte-1xieti7");
    			add_location(span0, file$2, 33, 34, 671);
    			attr_dev(div1, "class", "e-card-sub-title svelte-1xieti7");
    			add_location(div1, file$2, 33, 4, 641);
    			attr_dev(div2, "class", "e-card-header-caption");
    			add_location(div2, file$2, 31, 3, 549);
    			attr_dev(div3, "class", "e-card-header-image balance svelte-1xieti7");
    			add_location(div3, file$2, 35, 3, 752);
    			attr_dev(div4, "class", "e-card-header");
    			add_location(div4, file$2, 30, 2, 518);
    			attr_dev(span1, "class", "svelte-1xieti7");
    			add_location(span1, file$2, 38, 10, 858);
    			attr_dev(span2, "v-html", "this.duOSys.status.balance.unconfirmed");
    			attr_dev(span2, "class", "svelte-1xieti7");
    			add_location(span2, file$2, 38, 40, 888);
    			add_location(strong0, file$2, 38, 32, 880);
    			add_location(small0, file$2, 38, 3, 851);
    			attr_dev(span3, "class", "svelte-1xieti7");
    			add_location(span3, file$2, 39, 10, 977);
    			attr_dev(span4, "v-html", "this.duOSys.status.txsnumber");
    			attr_dev(span4, "class", "svelte-1xieti7");
    			add_location(span4, file$2, 39, 45, 1012);
    			add_location(strong1, file$2, 39, 37, 1004);
    			add_location(small1, file$2, 39, 3, 970);
    			attr_dev(div5, "class", "flx flc e-card-content svelte-1xieti7");
    			add_location(div5, file$2, 37, 2, 811);
    			attr_dev(div6, "class", "e-card flx flc justifyBetween duoCard svelte-1xieti7");
    			add_location(div6, file$2, 29, 1, 464);
    			attr_dev(div7, "class", "rwrap flx");
    			add_location(div7, file$2, 28, 0, 439);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div7, anchor);
    			append_dev(div7, div6);
    			append_dev(div6, div4);
    			append_dev(div4, div2);
    			append_dev(div2, div0);
    			append_dev(div2, t1);
    			append_dev(div2, div1);
    			append_dev(div1, span0);
    			append_dev(div1, t2);
    			append_dev(div4, t3);
    			append_dev(div4, div3);
    			append_dev(div6, t4);
    			append_dev(div6, div5);
    			append_dev(div5, small0);
    			append_dev(small0, span1);
    			append_dev(small0, strong0);
    			append_dev(strong0, span2);
    			append_dev(div5, t6);
    			append_dev(div5, small1);
    			append_dev(small1, span3);
    			append_dev(small1, strong1);
    			append_dev(strong1, span4);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div7);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$2.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelBalance extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$2, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelBalance",
    			options,
    			id: create_fragment$2.name
    		});
    	}
    }

    /* src/components/panels/PanelSend.svelte generated by Svelte v3.14.1 */
    const file$3 = "src/components/panels/PanelSend.svelte";

    // (12:3) <Button id="dialogbtn" v-on:click.native="btnClick">
    function create_default_slot(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("Send");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot.name,
    		type: "slot",
    		source: "(12:3) <Button id=\\\"dialogbtn\\\" v-on:click.native=\\\"btnClick\\\">",
    		ctx
    	});

    	return block;
    }

    function create_fragment$3(ctx) {
    	let div2;
    	let div1;
    	let t0;
    	let div0;
    	let t1;
    	let current;

    	const textfield0 = new smelte.TextField({
    			props: {
    				id: "sendAddress",
    				label: "Test label",
    				placeholder: "Enter DUO address",
    				outlined: true,
    				hint: "Test hint",
    				class: "e-outline flx noMargin fii bgfff"
    			},
    			$$inline: true
    		});

    	const textfield1 = new smelte.TextField({
    			props: {
    				label: "Test label",
    				class: "e-outline noMargin fii bgfff"
    			},
    			$$inline: true
    		});

    	const button = new smelte.Button({
    			props: {
    				id: "dialogbtn",
    				"v-on:click.native": "btnClick",
    				$$slots: { default: [create_default_slot] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			div2 = element("div");
    			div1 = element("div");
    			create_component(textfield0.$$.fragment);
    			t0 = space();
    			div0 = element("div");
    			create_component(textfield1.$$.fragment);
    			t1 = space();
    			create_component(button.$$.fragment);
    			attr_dev(div0, "class", "flx fii");
    			add_location(div0, file$3, 9, 3, 314);
    			attr_dev(div1, "class", "flx flc fii justifyBetween");
    			add_location(div1, file$3, 7, 1, 119);
    			attr_dev(div2, "class", "rwrap flx");
    			add_location(div2, file$3, 6, 0, 94);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div2, anchor);
    			append_dev(div2, div1);
    			mount_component(textfield0, div1, null);
    			append_dev(div1, t0);
    			append_dev(div1, div0);
    			mount_component(textfield1, div0, null);
    			append_dev(div0, t1);
    			mount_component(button, div0, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const button_changes = {};

    			if (changed.$$scope) {
    				button_changes.$$scope = { changed, ctx };
    			}

    			button.$set(button_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(textfield0.$$.fragment, local);
    			transition_in(textfield1.$$.fragment, local);
    			transition_in(button.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(textfield0.$$.fragment, local);
    			transition_out(textfield1.$$.fragment, local);
    			transition_out(button.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div2);
    			destroy_component(textfield0);
    			destroy_component(textfield1);
    			destroy_component(button);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$3.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelSend extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$3, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelSend",
    			options,
    			id: create_fragment$3.name
    		});
    	}
    }

    /* src/components/panels/PanelNetworkHashrate.svelte generated by Svelte v3.14.1 */

    function create_fragment$4(ctx) {
    	const block = {
    		c: noop,
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: noop,
    		p: noop,
    		i: noop,
    		o: noop,
    		d: noop
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$4.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelNetworkHashrate extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$4, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelNetworkHashrate",
    			options,
    			id: create_fragment$4.name
    		});
    	}
    }

    /* src/components/panels/PanelLocalHashrate.svelte generated by Svelte v3.14.1 */

    function create_fragment$5(ctx) {
    	const block = {
    		c: noop,
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: noop,
    		p: noop,
    		i: noop,
    		o: noop,
    		d: noop
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$5.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PanelLocalHashrate extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$5, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelLocalHashrate",
    			options,
    			id: create_fragment$5.name
    		});
    	}
    }

    /* src/components/panels/PanelStatus.svelte generated by Svelte v3.14.1 */

    const file$4 = "src/components/panels/PanelStatus.svelte";

    function create_fragment$6(ctx) {
    	let div;
    	let ul;
    	let li0;
    	let span0;
    	let t0;
    	let span1;
    	let t2;
    	let strong0;
    	let t3_value = ctx.status.ver + "";
    	let t3;
    	let t4;
    	let li1;
    	let span2;
    	let t5;
    	let span3;
    	let t7;
    	let strong1;
    	let t8_value = ctx.status.walletver.podjsonrpcapi.versionstring + "";
    	let t8;
    	let t9;
    	let li2;
    	let span4;
    	let t10;
    	let span5;
    	let t12;
    	let strong2;
    	let t13_value = ctx.status.uptime + "";
    	let t13;
    	let t14;
    	let li3;
    	let span6;
    	let t15;
    	let span7;
    	let t17;
    	let strong3;
    	let t18_value = ctx.status.net + "";
    	let t18;
    	let t19;
    	let li4;
    	let span8;
    	let t20;
    	let span9;
    	let t22;
    	let strong4;
    	let t23_value = ctx.status.ver + "";
    	let t23;
    	let t24;
    	let li5;
    	let span10;
    	let t25;
    	let span11;
    	let t27;
    	let strong5;
    	let t28_value = ctx.status.walletver.podjsonrpcapi.versionstring + "";
    	let t28;
    	let t29;
    	let li6;
    	let span12;
    	let t30;
    	let span13;
    	let t32;
    	let strong6;
    	let t33_value = ctx.status.uptime + "";
    	let t33;
    	let t34;
    	let li7;
    	let span14;
    	let t35;
    	let span15;
    	let t37;
    	let strong7;
    	let t38_value = ctx.status.net + "";
    	let t38;

    	const block = {
    		c: function create() {
    			div = element("div");
    			ul = element("ul");
    			li0 = element("li");
    			span0 = element("span");
    			t0 = space();
    			span1 = element("span");
    			span1.textContent = "Version:";
    			t2 = space();
    			strong0 = element("strong");
    			t3 = text(t3_value);
    			t4 = space();
    			li1 = element("li");
    			span2 = element("span");
    			t5 = space();
    			span3 = element("span");
    			span3.textContent = "Wallet version:";
    			t7 = space();
    			strong1 = element("strong");
    			t8 = text(t8_value);
    			t9 = space();
    			li2 = element("li");
    			span4 = element("span");
    			t10 = space();
    			span5 = element("span");
    			span5.textContent = "Uptime:";
    			t12 = space();
    			strong2 = element("strong");
    			t13 = text(t13_value);
    			t14 = space();
    			li3 = element("li");
    			span6 = element("span");
    			t15 = space();
    			span7 = element("span");
    			span7.textContent = "Memory:";
    			t17 = space();
    			strong3 = element("strong");
    			t18 = text(t18_value);
    			t19 = space();
    			li4 = element("li");
    			span8 = element("span");
    			t20 = space();
    			span9 = element("span");
    			span9.textContent = "Disk:";
    			t22 = space();
    			strong4 = element("strong");
    			t23 = text(t23_value);
    			t24 = space();
    			li5 = element("li");
    			span10 = element("span");
    			t25 = space();
    			span11 = element("span");
    			span11.textContent = "Chain:";
    			t27 = space();
    			strong5 = element("strong");
    			t28 = text(t28_value);
    			t29 = space();
    			li6 = element("li");
    			span12 = element("span");
    			t30 = space();
    			span13 = element("span");
    			span13.textContent = "Blocks:";
    			t32 = space();
    			strong6 = element("strong");
    			t33 = text(t33_value);
    			t34 = space();
    			li7 = element("li");
    			span14 = element("span");
    			t35 = space();
    			span15 = element("span");
    			span15.textContent = "Connections:";
    			t37 = space();
    			strong7 = element("strong");
    			t38 = text(t38_value);
    			attr_dev(span0, "class", "rcx2");
    			add_location(span0, file$4, 7, 11, 180);
    			attr_dev(span1, "class", "rcx4");
    			add_location(span1, file$4, 8, 11, 218);
    			attr_dev(strong0, "class", "rcx6");
    			add_location(strong0, file$4, 9, 11, 265);
    			attr_dev(li0, "class", "flx fwd spb htg rr");
    			add_location(li0, file$4, 6, 2, 137);
    			attr_dev(span2, "class", "rcx2");
    			add_location(span2, file$4, 12, 11, 377);
    			attr_dev(span3, "class", "rcx4");
    			add_location(span3, file$4, 13, 11, 415);
    			attr_dev(strong1, "class", "rcx6");
    			add_location(strong1, file$4, 14, 11, 469);
    			attr_dev(li1, "class", "flx fwd spb htg rr");
    			add_location(li1, file$4, 11, 10, 334);
    			attr_dev(span4, "class", "rcx2");
    			add_location(span4, file$4, 17, 11, 615);
    			attr_dev(span5, "class", "rcx4");
    			add_location(span5, file$4, 18, 11, 653);
    			attr_dev(strong2, "class", "rcx6");
    			add_location(strong2, file$4, 19, 11, 699);
    			attr_dev(li2, "class", "flx fwd spb htg rr");
    			add_location(li2, file$4, 16, 10, 572);
    			attr_dev(span6, "class", "rcx2");
    			add_location(span6, file$4, 22, 11, 814);
    			attr_dev(span7, "class", "rcx4");
    			add_location(span7, file$4, 23, 11, 852);
    			attr_dev(strong3, "class", "rcx6");
    			add_location(strong3, file$4, 24, 11, 898);
    			attr_dev(li3, "class", "flx fwd spb htg rr");
    			add_location(li3, file$4, 21, 10, 771);
    			attr_dev(span8, "class", "rcx2");
    			add_location(span8, file$4, 28, 19, 1019);
    			attr_dev(span9, "class", "rcx4");
    			add_location(span9, file$4, 29, 19, 1065);
    			attr_dev(strong4, "class", "rcx6");
    			add_location(strong4, file$4, 30, 19, 1117);
    			attr_dev(li4, "class", "flx fwd spb htg rr");
    			add_location(li4, file$4, 27, 10, 968);
    			attr_dev(span10, "class", "rcx2");
    			add_location(span10, file$4, 33, 19, 1253);
    			attr_dev(span11, "class", "rcx4");
    			add_location(span11, file$4, 34, 19, 1299);
    			attr_dev(strong5, "class", "rcx6");
    			add_location(strong5, file$4, 35, 19, 1352);
    			attr_dev(li5, "class", "flx fwd spb htg rr");
    			add_location(li5, file$4, 32, 18, 1202);
    			attr_dev(span12, "class", "rcx2");
    			add_location(span12, file$4, 38, 19, 1522);
    			attr_dev(span13, "class", "rcx4");
    			add_location(span13, file$4, 39, 19, 1568);
    			attr_dev(strong6, "class", "rcx6");
    			add_location(strong6, file$4, 40, 19, 1622);
    			attr_dev(li6, "class", "flx fwd spb htg rr");
    			add_location(li6, file$4, 37, 18, 1471);
    			attr_dev(span14, "class", "rcx2");
    			add_location(span14, file$4, 43, 19, 1761);
    			attr_dev(span15, "class", "rcx4");
    			add_location(span15, file$4, 44, 19, 1807);
    			attr_dev(strong7, "class", "rcx6");
    			add_location(strong7, file$4, 45, 19, 1866);
    			attr_dev(li7, "class", "flx fwd spb htg rr");
    			add_location(li7, file$4, 42, 18, 1710);
    			attr_dev(ul, "class", "rf flx flc noMargin noPadding justifyEvenly");
    			add_location(ul, file$4, 5, 1, 78);
    			attr_dev(div, "id", "panelstatus");
    			attr_dev(div, "class", "Info");
    			add_location(div, file$4, 4, 0, 40);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			append_dev(div, ul);
    			append_dev(ul, li0);
    			append_dev(li0, span0);
    			append_dev(li0, t0);
    			append_dev(li0, span1);
    			append_dev(li0, t2);
    			append_dev(li0, strong0);
    			append_dev(strong0, t3);
    			append_dev(ul, t4);
    			append_dev(ul, li1);
    			append_dev(li1, span2);
    			append_dev(li1, t5);
    			append_dev(li1, span3);
    			append_dev(li1, t7);
    			append_dev(li1, strong1);
    			append_dev(strong1, t8);
    			append_dev(ul, t9);
    			append_dev(ul, li2);
    			append_dev(li2, span4);
    			append_dev(li2, t10);
    			append_dev(li2, span5);
    			append_dev(li2, t12);
    			append_dev(li2, strong2);
    			append_dev(strong2, t13);
    			append_dev(ul, t14);
    			append_dev(ul, li3);
    			append_dev(li3, span6);
    			append_dev(li3, t15);
    			append_dev(li3, span7);
    			append_dev(li3, t17);
    			append_dev(li3, strong3);
    			append_dev(strong3, t18);
    			append_dev(ul, t19);
    			append_dev(ul, li4);
    			append_dev(li4, span8);
    			append_dev(li4, t20);
    			append_dev(li4, span9);
    			append_dev(li4, t22);
    			append_dev(li4, strong4);
    			append_dev(strong4, t23);
    			append_dev(ul, t24);
    			append_dev(ul, li5);
    			append_dev(li5, span10);
    			append_dev(li5, t25);
    			append_dev(li5, span11);
    			append_dev(li5, t27);
    			append_dev(li5, strong5);
    			append_dev(strong5, t28);
    			append_dev(ul, t29);
    			append_dev(ul, li6);
    			append_dev(li6, span12);
    			append_dev(li6, t30);
    			append_dev(li6, span13);
    			append_dev(li6, t32);
    			append_dev(li6, strong6);
    			append_dev(strong6, t33);
    			append_dev(ul, t34);
    			append_dev(ul, li7);
    			append_dev(li7, span14);
    			append_dev(li7, t35);
    			append_dev(li7, span15);
    			append_dev(li7, t37);
    			append_dev(li7, strong7);
    			append_dev(strong7, t38);
    		},
    		p: function update(changed, ctx) {
    			if (changed.status && t3_value !== (t3_value = ctx.status.ver + "")) set_data_dev(t3, t3_value);
    			if (changed.status && t8_value !== (t8_value = ctx.status.walletver.podjsonrpcapi.versionstring + "")) set_data_dev(t8, t8_value);
    			if (changed.status && t13_value !== (t13_value = ctx.status.uptime + "")) set_data_dev(t13, t13_value);
    			if (changed.status && t18_value !== (t18_value = ctx.status.net + "")) set_data_dev(t18, t18_value);
    			if (changed.status && t23_value !== (t23_value = ctx.status.ver + "")) set_data_dev(t23, t23_value);
    			if (changed.status && t28_value !== (t28_value = ctx.status.walletver.podjsonrpcapi.versionstring + "")) set_data_dev(t28, t28_value);
    			if (changed.status && t33_value !== (t33_value = ctx.status.uptime + "")) set_data_dev(t33, t33_value);
    			if (changed.status && t38_value !== (t38_value = ctx.status.net + "")) set_data_dev(t38, t38_value);
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$6.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$1($$self, $$props, $$invalidate) {
    	let { status } = $$props;
    	const writable_props = ["status"];

    	Object.keys($$props).forEach(key => {
    		if (!~writable_props.indexOf(key) && key.slice(0, 2) !== "$$") console.warn(`<PanelStatus> was created with unknown prop '${key}'`);
    	});

    	$$self.$set = $$props => {
    		if ("status" in $$props) $$invalidate("status", status = $$props.status);
    	};

    	$$self.$capture_state = () => {
    		return { status };
    	};

    	$$self.$inject_state = $$props => {
    		if ("status" in $$props) $$invalidate("status", status = $$props.status);
    	};

    	return { status };
    }

    class PanelStatus extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$1, create_fragment$6, safe_not_equal, { status: 0 });

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelStatus",
    			options,
    			id: create_fragment$6.name
    		});

    		const { ctx } = this.$$;
    		const props = options.props || ({});

    		if (ctx.status === undefined && !("status" in props)) {
    			console.warn("<PanelStatus> was created without expected prop 'status'");
    		}
    	}

    	get status() {
    		throw new Error("<PanelStatus>: Props cannot be read directly from the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}

    	set status(value) {
    		throw new Error("<PanelStatus>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    const loading = writable(false);

    const lasTxs = readable([], function start(set) {
        const interval = setInterval(() => {
            fetch(`http://127.0.0.1:3999/lastxs`)
                .then(resp => resp.json())
                .then(data => {
                    set(data);
                });
        }, 100);
        return function stop() {
            clearInterval(interval);
        };
    });

    /* src/components/panels/PanelLatestTx.svelte generated by Svelte v3.14.1 */
    const file$5 = "src/components/panels/PanelLatestTx.svelte";

    function create_fragment$7(ctx) {
    	let div;
    	let current;

    	const datatable = new smelte.DataTable({
    			props: {
    				lasTxs,
    				loading,
    				columns: [
    					{
    						label: "ID",
    						field: "id",
    						class: "md:w-10"
    					},
    					{
    						label: "Ep.",
    						value: func,
    						class: "md:w-10",
    						editable: false
    					},
    					{ field: "name", class: "md:w-10" },
    					{
    						field: "summary",
    						textarea: true,
    						value: func_1,
    						class: "text-sm text-gray-700 caption md:w-full sm:w-64"
    					},
    					{
    						field: "thumbnail",
    						value: func_2,
    						class: "w-48",
    						sortable: false,
    						editable: false
    					}
    				]
    			},
    			$$inline: true
    		});

    	datatable.$on("update", ctx.update_handler);

    	const block = {
    		c: function create() {
    			div = element("div");
    			create_component(datatable.$$.fragment);
    			attr_dev(div, "class", "rwrap flx");
    			add_location(div, file$5, 6, 0, 157);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div, anchor);
    			mount_component(datatable, div, null);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const datatable_changes = {};
    			if (changed.lasTxs) datatable_changes.lasTxs = lasTxs;
    			datatable.$set(datatable_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(datatable.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(datatable.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div);
    			destroy_component(datatable);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$7.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    const func = v => `S${v.season}E${v.number}`;
    const func_1 = v => v && v.summary ? v.summary : "";

    const func_2 = v => v && v.image
    ? `<img src="${v.image.medium.replace("http", "https")}" height="70" alt="${v.name}">`
    : "";

    function instance$2($$self) {
    	const update_handler = ({ detail }) => {
    		const { column, item, value } = detail;
    		const index = lasTxs.findIndex(i => i.id === item.id);
    		lasTxs[index][column.field] = value;
    	};

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		
    	};

    	return { update_handler };
    }

    class PanelLatestTx extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$2, create_fragment$7, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelLatestTx",
    			options,
    			id: create_fragment$7.name
    		});
    	}
    }

    const status = readable([], function start(set) {
        const interval = setInterval(() => {
            fetch(`http://127.0.0.1:3999/status`)
                .then(resp => resp.json())
                .then(data => {
                    set(data);
                });
        }, 1000);
        return function stop() {
            clearInterval(interval);
        };
    });

    /* src/components/pages/PageOverview.svelte generated by Svelte v3.14.1 */
    const file$6 = "src/components/pages/PageOverview.svelte";

    function create_fragment$8(ctx) {
    	let t0;
    	let main;
    	let div0;
    	let t1;
    	let div1;
    	let t2;
    	let div2;
    	let t3;
    	let div3;
    	let t4;
    	let div4;
    	let t5;
    	let div5;
    	let t6;
    	let div6;
    	let t7;
    	let div7;
    	let current;

    	const panelbalance = new PanelBalance({
    			props: { balance: status.balance },
    			$$inline: true
    		});

    	const panelsend = new PanelSend({ $$inline: true });

    	const panelnetworkhashrate = new PanelNetworkHashrate({
    			props: { hashrate: status.networkhashrate },
    			$$inline: true
    		});

    	const panellocalhashrate = new PanelLocalHashrate({
    			props: { hashrate: status.networkhashrate },
    			$$inline: true
    		});

    	const panelstatus = new PanelStatus({ props: { status }, $$inline: true });
    	const panellatesttx = new PanelLatestTx({ $$inline: true });

    	const block = {
    		c: function create() {
    			t0 = space();
    			main = element("main");
    			div0 = element("div");
    			create_component(panelbalance.$$.fragment);
    			t1 = space();
    			div1 = element("div");
    			create_component(panelsend.$$.fragment);
    			t2 = space();
    			div2 = element("div");
    			create_component(panelnetworkhashrate.$$.fragment);
    			t3 = space();
    			div3 = element("div");
    			create_component(panellocalhashrate.$$.fragment);
    			t4 = space();
    			div4 = element("div");
    			create_component(panelstatus.$$.fragment);
    			t5 = space();
    			div5 = element("div");
    			create_component(panellatesttx.$$.fragment);
    			t6 = space();
    			div6 = element("div");
    			t7 = space();
    			div7 = element("div");
    			document.title = "Overview";
    			attr_dev(div0, "id", "panelwalletstatus");
    			attr_dev(div0, "class", "Balance");
    			add_location(div0, file$6, 25, 8, 591);
    			attr_dev(div1, "id", "panelsend");
    			attr_dev(div1, "class", "Send");
    			add_location(div1, file$6, 28, 8, 712);
    			attr_dev(div2, "id", "panelnetworkhashrate");
    			attr_dev(div2, "class", "NetHash");
    			add_location(div2, file$6, 31, 8, 796);
    			attr_dev(div3, "id", "panellocalhashrate");
    			attr_dev(div3, "class", "LocalHash");
    			add_location(div3, file$6, 34, 8, 938);
    			attr_dev(div4, "id", "panelstatus");
    			attr_dev(div4, "class", "Status");
    			add_location(div4, file$6, 37, 8, 1078);
    			attr_dev(div5, "id", "paneltxsex");
    			attr_dev(div5, "class", "Txs");
    			add_location(div5, file$6, 40, 8, 1183);
    			attr_dev(div6, "class", "Info");
    			add_location(div6, file$6, 43, 8, 1270);
    			attr_dev(div7, "class", "Time");
    			add_location(div7, file$6, 45, 8, 1312);
    			attr_dev(main, "class", "pageOverview");
    			add_location(main, file$6, 22, 8, 553);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t0, anchor);
    			insert_dev(target, main, anchor);
    			append_dev(main, div0);
    			mount_component(panelbalance, div0, null);
    			append_dev(main, t1);
    			append_dev(main, div1);
    			mount_component(panelsend, div1, null);
    			append_dev(main, t2);
    			append_dev(main, div2);
    			mount_component(panelnetworkhashrate, div2, null);
    			append_dev(main, t3);
    			append_dev(main, div3);
    			mount_component(panellocalhashrate, div3, null);
    			append_dev(main, t4);
    			append_dev(main, div4);
    			mount_component(panelstatus, div4, null);
    			append_dev(main, t5);
    			append_dev(main, div5);
    			mount_component(panellatesttx, div5, null);
    			append_dev(main, t6);
    			append_dev(main, div6);
    			append_dev(main, t7);
    			append_dev(main, div7);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(panelbalance.$$.fragment, local);
    			transition_in(panelsend.$$.fragment, local);
    			transition_in(panelnetworkhashrate.$$.fragment, local);
    			transition_in(panellocalhashrate.$$.fragment, local);
    			transition_in(panelstatus.$$.fragment, local);
    			transition_in(panellatesttx.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(panelbalance.$$.fragment, local);
    			transition_out(panelsend.$$.fragment, local);
    			transition_out(panelnetworkhashrate.$$.fragment, local);
    			transition_out(panellocalhashrate.$$.fragment, local);
    			transition_out(panelstatus.$$.fragment, local);
    			transition_out(panellatesttx.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t0);
    			if (detaching) detach_dev(main);
    			destroy_component(panelbalance);
    			destroy_component(panelsend);
    			destroy_component(panelnetworkhashrate);
    			destroy_component(panellocalhashrate);
    			destroy_component(panelstatus);
    			destroy_component(panellatesttx);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$8.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageOverview extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$8, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageOverview",
    			options,
    			id: create_fragment$8.name
    		});
    	}
    }

    /* src/components/panels/PanelTxs.svelte generated by Svelte v3.14.1 */

    const file$7 = "src/components/panels/PanelTxs.svelte";

    function create_fragment$9(ctx) {
    	let div1;
    	let div0;

    	const block = {
    		c: function create() {
    			div1 = element("div");
    			div0 = element("div");
    			attr_dev(div0, "id", "txs");
    			add_location(div0, file$7, 30, 23, 1044);
    			attr_dev(div1, "class", "rwrap flx");
    			add_location(div1, file$7, 30, 0, 1021);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div1, anchor);
    			append_dev(div1, div0);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div1);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$9.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$3($$self, $$props, $$invalidate) {
    	const txs = new Vue({
    			el: "#txs",
    			name: "PanelTransactions",
    			data() {
    				return {
    					pageSettings: {
    						pageSize: 10,
    						pageSizes: [10, 20, 50, 100],
    						pageCount: 5
    					},
    					ddldata: ["All", "generated", "sent", "received", "immature"]
    				};
    			},
    			template: `<div class="rwrap">
	<div class="select-wrap">
		<ejs-dropdownlist id='ddlelement' :dataSource='ddldata' placeholder='Select category to filter'></ejs-dropdownlist>
    </div>
	<ejs-grid :dataSource="this.transactions.txsEx.txs" height="100%" :allowPaging="true" :pageSettings='pageSettings'>
		<e-columns>
			<e-column field='category' headerText='Category' textAlign='Right' width=90></e-column>
			<e-column field='time' headerText='Time' format='unix'  textAlign='Right' width=90></e-column>
			<e-column field='txid' headerText='TxID' textAlign='Right' width=240></e-column>
			<e-column field='amount' headerText='Amount' textAlign='Right' width=90></e-column>
		</e-columns>
	</ejs-grid>
	</div>`,
    			props: { transactions: Object }
    		});

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		
    	};

    	return { txs };
    }

    class PanelTxs extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$3, create_fragment$9, safe_not_equal, { txs: 0 });

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelTxs",
    			options,
    			id: create_fragment$9.name
    		});
    	}

    	get txs() {
    		return this.$$.ctx.txs;
    	}

    	set txs(value) {
    		throw new Error("<PanelTxs>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* src/components/pages/PageTransactions.svelte generated by Svelte v3.14.1 */

    function create_fragment$a(ctx) {
    	let t;
    	let current;
    	const paneltxs = new PanelTxs({ $$inline: true });

    	const block = {
    		c: function create() {
    			t = space();
    			create_component(paneltxs.$$.fragment);
    			document.title = "Transactions";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    			mount_component(paneltxs, target, anchor);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(paneltxs.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(paneltxs.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    			destroy_component(paneltxs, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$a.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageTransactions extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$a, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageTransactions",
    			options,
    			id: create_fragment$a.name
    		});
    	}
    }

    /* src/components/pages/PageAddressBook.svelte generated by Svelte v3.14.1 */

    function create_fragment$b(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("\n\nPage PageAddressBook");
    			document.title = "Address Book";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$b.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageAddressBook extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$b, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageAddressBook",
    			options,
    			id: create_fragment$b.name
    		});
    	}
    }

    /* src/components/pages/PageExplorer.svelte generated by Svelte v3.14.1 */

    function create_fragment$c(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("\n\n\nPage PageExplorer");
    			document.title = "Explorer";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$c.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageExplorer extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$c, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageExplorer",
    			options,
    			id: create_fragment$c.name
    		});
    	}
    }

    /* src/components/panels/PanelSettings.svelte generated by Svelte v3.14.1 */

    const file$8 = "src/components/panels/PanelSettings.svelte";

    function create_fragment$d(ctx) {
    	let div1;
    	let div0;

    	const block = {
    		c: function create() {
    			div1 = element("div");
    			div0 = element("div");
    			attr_dev(div0, "id", "sts");
    			add_location(div0, file$8, 49, 23, 4135);
    			attr_dev(div1, "class", "rwrap flx");
    			add_location(div1, file$8, 49, 0, 4112);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, div1, anchor);
    			append_dev(div1, div0);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(div1);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$d.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$4($$self, $$props, $$invalidate) {
    	const sts = new Vue({
    			el: "#sts",
    			name: "PanelSettings",
    			data() {
    				return { duoSystem };
    			},
    			created() {
    				
    			},
    			methods: {},
    			template: `<div class="rwrap">asasasas
<div v-html="this.duoSystem.config.daemon.schema"></div>
 <svelte-form-generator class="flx flc fii" :schema="this.duoSystem.config.daemon.schema" :model="this.duoSystem.config.daemon.config"></svelte-form-generator>
		</div>`,
    			style: `
button {
  margin: 25px 5px 20px 20px;
}
@font-face {
    font-family: 'btn-icon';
    src:
    url(data:application/x-font-ttf;charset=utf-8;base64,AAEAAAAKAIAAAwAgT1MvMj1tSfgAAAEoAAAAVmNtYXDnH+dzAAABoAAAAEJnbHlm1v48pAAAAfgAAAQYaGVhZBOPfZcAAADQAAAANmhoZWEIUQQJAAAArAAAACRobXR4IAAAAAAAAYAAAAAgbG9jYQN6ApQAAAHkAAAAEm1heHABFQCqAAABCAAAACBuYW1l07lFxAAABhAAAAIxcG9zdK9uovoAAAhEAAAAgAABAAAEAAAAAFwEAAAAAAAD9AABAAAAAAAAAAAAAAAAAAAACAABAAAAAQAAJ1LUzF8PPPUACwQAAAAAANg+nFMAAAAA2D6cUwAAAAAD9AP0AAAACAACAAAAAAAAAAEAAAAIAJ4AAwAAAAAAAgAAAAoACgAAAP8AAAAAAAAAAQQAAZAABQAAAokCzAAAAI8CiQLMAAAB6wAyAQgAAAIABQMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAUGZFZABA5wDnBgQAAAAAXAQAAAAAAAABAAAAAAAABAAAAAQAAAAEAAAABAAAAAQAAAAEAAAABAAAAAQAAAAAAAACAAAAAwAAABQAAwABAAAAFAAEAC4AAAAEAAQAAQAA5wb//wAA5wD//wAAAAEABAAAAAEAAgADAAQABQAGAAcAAAAAAAAADgAkADIAhAEuAewCDAAAAAEAAAAAA2ED9AACAAA3CQGeAsT9PAwB9AH0AAACAAAAAAPHA/QAAwAHAAAlIREhASERIQJpAV7+ov3QAV7+ogwD6PwYA+gAAAEAAAAAA4sD9AACAAATARF0AxgCAP4MA+gAAAABAAAAAAP0A/QAQwAAExEfDyE/DxEvDyEPDgwBAgMFBQcICQkLCwwMDQ4NAtoNDg0MDAsLCQkIBwUFAwIBAQIDBQUHCAkJCwsMDA0ODf0mDQ4NDAwLCwkJCAcFBQMCA239Jg4NDQ0LCwsJCQgHBQUDAgEBAgMFBQcICQkLCwsNDQ0OAtoODQ0NCwsLCQkIBwUFAwIBAQIDBQUHCAkJCwsLDQ0NAAIAAAAAA/MDxQADAIwAADczESMBDwMVFw8METM3HwQ3Fz8KPQEvBT8LLwg3NT8INS8FNT8NNS8JByU/BDUvCyMPAQytrQH5AgoEAQEBARghERESEyIJCSgQBiEHNQceOZPbDgUICw0LCQUDBAICBAkGAgEBAQMOBAkIBgcDAwEBAQEDAwMJAgEBAxYLBQQEAwMCAgIEBAoBAQEECgcHBgUFBAMDAQEBAQQFBwkFBQUGEf6tDwkEAwIBAQMDCgwVAwcGDAsNBwdaAYcB3gEFAwN2HwoELDodGxwaLwkIGwz+igEBHwMBAQECAQEDBgoKDAYICAgFCAkICwUEBAQFAwYDBwgIDAgHCAcGBgYFBQkEAgYCBAwJBgUGBwkJCgkICAcLBAIFAwIEBAQFBQcGBwgHBgYGBgoJCAYCAgEBAQFGMRkaGw0NDA0LIh4xBAQCBAEBAgADAAAAAAOKA/MAHABCAJ0AAAEzHwIRDwMhLwIDNzM/CjUTHwcVIwcVIy8HETcXMz8KNScxBxEfDjsBHQEfDTMhMz8OES8PIz0BLw4hA0EDBQQDAQIEBf5eBQQCAW4RDg0LCQgGBQUDBAFeBAMDAwIBAQGL7Y0EAwQCAgIBAYYKChEQDQsJCAcEBAUCYt8BAQIDBAUFBQcHBwgICQgKjQECAgMEBAUFBgYHBgcIBwGcCAcHBwYGBgUFBAQDAgIBAQEBAgIDBAQFBQYGBgcHBwgmAQMDAwUFBgYHBwgICQkJ/tQCiwMEBf3XAwYEAgIEBgFoAQEDBQYGBwgIBw0KhQEiAQEBAgMDAwTV+94BAQECAwMDBAGyAQECBAYHCAgJCgkQCaQC6/47CQkICQcIBwYGBQQEAwICUAgHBwcGBgYFBQQEAwMBAgIBAwMEBAUFBQcGBwcHCAImCAcHBwYGBgUFBAQDAgIBAdUJCQgICAgGBwYFBAQDAgEBAAAAAAIAAAAAA6cD9AADAAwAADchNSElAQcJAScBESNZA078sgGB/uMuAXkBgDb+1EwMTZcBCD3+ngFiPf7pAxMAAAAAABIA3gABAAAAAAAAAAEAAAABAAAAAAABAAgAAQABAAAAAAACAAcACQABAAAAAAADAAgAEAABAAAAAAAEAAgAGAABAAAAAAAFAAsAIAABAAAAAAAGAAgAKwABAAAAAAAKACwAMwABAAAAAAALABIAXwADAAEECQAAAAIAcQADAAEECQABABAAcwADAAEECQACAA4AgwADAAEECQADABAAkQADAAEECQAEABAAoQADAAEECQAFABYAsQADAAEECQAGABAAxwADAAEECQAKAFgA1wADAAEECQALACQBLyBidG4taWNvblJlZ3VsYXJidG4taWNvbmJ0bi1pY29uVmVyc2lvbiAxLjBidG4taWNvbkZvbnQgZ2VuZXJhdGVkIHVzaW5nIFN5bmNmdXNpb24gTWV0cm8gU3R1ZGlvd3d3LnN5bmNmdXNpb24uY29tACAAYgB0AG4ALQBpAGMAbwBuAFIAZQBnAHUAbABhAHIAYgB0AG4ALQBpAGMAbwBuAGIAdABuAC0AaQBjAG8AbgBWAGUAcgBzAGkAbwBuACAAMQAuADAAYgB0AG4ALQBpAGMAbwBuAEYAbwBuAHQAIABnAGUAbgBlAHIAYQB0AGUAZAAgAHUAcwBpAG4AZwAgAFMAeQBuAGMAZgB1AHMAaQBvAG4AIABNAGUAdAByAG8AIABTAHQAdQBkAGkAbwB3AHcAdwAuAHMAeQBuAGMAZgB1AHMAaQBvAG4ALgBjAG8AbQAAAAACAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBAgEDAQQBBQEGAQcBCAEJAAptZWRpYS1wbGF5C21lZGlhLXBhdXNlDmFycm93aGVhZC1sZWZ0BHN0b3AJbGlrZS0tLTAxBGNvcHkQLWRvd25sb2FkLTAyLXdmLQAA) format('truetype');
    font-weight: normal;
    font-style: normal;
}
.e-btn-sb-icon {
    font-family: 'btn-icon' !important;
    speak: none;
    font-size: 55px;
    font-style: normal;
    font-weight: normal;
    font-variant: normal;
    text-transform: none;
    line-height: 1;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}
/*For Right Icon Button*/
.e-send-icon::before {
  content: '\eb35';
}
/*For Left Icon Button*/
.e-receive-icon::before {
  content: '\eb0e';
}`
    		});

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		
    	};

    	return { sts };
    }

    class PanelSettings extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$4, create_fragment$d, safe_not_equal, { sts: 0 });

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PanelSettings",
    			options,
    			id: create_fragment$d.name
    		});
    	}

    	get sts() {
    		return this.$$.ctx.sts;
    	}

    	set sts(value) {
    		throw new Error("<PanelSettings>: Props cannot be set directly on the component instance unless compiling with 'accessors: true' or '<svelte:options accessors/>'");
    	}
    }

    /* src/components/pages/PageSettings.svelte generated by Svelte v3.14.1 */

    function create_fragment$e(ctx) {
    	let t;
    	let current;
    	const panelsettings = new PanelSettings({ $$inline: true });

    	const block = {
    		c: function create() {
    			t = space();
    			create_component(panelsettings.$$.fragment);
    			document.title = "Settings";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    			mount_component(panelsettings, target, anchor);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(panelsettings.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(panelsettings.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    			destroy_component(panelsettings, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$e.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageSettings extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$e, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageSettings",
    			options,
    			id: create_fragment$e.name
    		});
    	}
    }

    /* src/components/pages/PageNotFound.svelte generated by Svelte v3.14.1 */
    const file$9 = "src/components/pages/PageNotFound.svelte";

    // (9:0) <Button>
    function create_default_slot$1(ctx) {
    	let t;

    	const block = {
    		c: function create() {
    			t = text("Go Home");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t, anchor);
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_default_slot$1.name,
    		type: "slot",
    		source: "(9:0) <Button>",
    		ctx
    	});

    	return block;
    }

    function create_fragment$f(ctx) {
    	let h1;
    	let t1;
    	let t2;
    	let current;

    	const button = new smelte.Button({
    			props: {
    				$$slots: { default: [create_default_slot$1] },
    				$$scope: { ctx }
    			},
    			$$inline: true
    		});

    	const block = {
    		c: function create() {
    			h1 = element("h1");
    			h1.textContent = "404 Not Found";
    			t1 = space();
    			create_component(button.$$.fragment);
    			t2 = space();
    			add_location(h1, file$9, 7, 0, 58);
    			document.title = "404";
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, h1, anchor);
    			insert_dev(target, t1, anchor);
    			mount_component(button, target, anchor);
    			insert_dev(target, t2, anchor);
    			current = true;
    		},
    		p: function update(changed, ctx) {
    			const button_changes = {};

    			if (changed.$$scope) {
    				button_changes.$$scope = { changed, ctx };
    			}

    			button.$set(button_changes);
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(button.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(button.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(h1);
    			if (detaching) detach_dev(t1);
    			destroy_component(button, detaching);
    			if (detaching) detach_dev(t2);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$f.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class PageNotFound extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$f, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "PageNotFound",
    			options,
    			id: create_fragment$f.name
    		});
    	}
    }

    function createPages() {
        const {subscribe, set} = writable(PageOverview);

        return {
            subscribe,
            overview: () => set(PageOverview),
            transactions: () => set(PageTransactions),
            addressbook: () => set(PageAddressBook),
            explorer: () => set(PageExplorer),
            settings: () => set(PageSettings),
            notfound: () => set(PageNotFound),
        };
    }

    const isPage = createPages();

    /* src/components/ico/Logo.svelte generated by Svelte v3.14.1 */

    const file$a = "src/components/ico/Logo.svelte";

    function create_fragment$g(ctx) {
    	let svg;
    	let path;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path = svg_element("path");
    			attr_dev(path, "class", "logofill");
    			attr_dev(path, "d", "M77.08,2.55c3.87,1.03 6.96,2.58 10.32,4.64c5.93,3.87 10.58,8.51 14.19,14.71c3.87,6.19 5.42,13.16 5.42,20.64c0,7.22 -1.81,14.18 -5.41,20.37c-3.61,6.45 -8.25,11.35 -14.19,14.96c-3.35,2.06 -6.96,3.87 -10.32,4.9c-3.87,1.03 -7.74,1.55 -11.61,1.55v-14.45c6.96,-0.26 13.42,-2.58 19.09,-8c5.67,-5.42 8.51,-11.87 8.51,-19.61c0,-7.74 -2.58,-14.19 -7.99,-19.6c-5.42,-5.42 -11.86,-8 -19.6,-8c-7.74,0 -14.44,2.58 -19.6,8c-5.42,5.42 -8,11.87 -8,19.6l0,85.9c-3.1,-3.1 -7.99,-7.74 -13.93,-13.67v-72.23c0,-3.87 0.52,-7.73 1.55,-11.35c1.03,-3.87 2.58,-7.22 4.64,-10.32c3.87,-5.93 8.52,-10.58 14.71,-14.45c6.19,-3.61 13.16,-5.16 20.64,-5.16c3.87,0 8,0.52 11.61,1.55zM78.37,42.28c0,7.22 -5.93,13.16 -13.15,13.16c-7.48,0.26 -13.16,-5.68 -13.16,-13.16c0,-7.22 5.94,-13.16 13.16,-13.16c7.22,0 13.15,5.93 13.15,13.16zM13.63,37.12l0,69.39c-6.19,-6.19 -11.09,-10.83 -13.93,-13.93l0,-55.46z");
    			add_location(path, file$a, 0, 109, 109);
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "id", "parallelCoinLogo");
    			attr_dev(svg, "viewBox", "0 0 108 128");
    			attr_dev(svg, "width", "108");
    			attr_dev(svg, "height", "128");
    			attr_dev(svg, "class", "svelte-15pj3m3");
    			add_location(svg, file$a, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$g.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Logo extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$g, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Logo",
    			options,
    			id: create_fragment$g.name
    		});
    	}
    }

    /* src/components/layout/Header.svelte generated by Svelte v3.14.1 */

    const file$b = "src/components/layout/Header.svelte";

    function create_fragment$h(ctx) {
    	let header;
    	let div0;
    	let t0;
    	let div1;
    	let t1;
    	let div4;
    	let div3;
    	let div2;
    	let t3;
    	let div5;
    	let h1;
    	let t5;
    	let div6;
    	let t6;
    	let div7;
    	let t7;
    	let div8;
    	let t8;
    	let div9;
    	let button;

    	const block = {
    		c: function create() {
    			header = element("header");
    			div0 = element("div");
    			t0 = space();
    			div1 = element("div");
    			t1 = space();
    			div4 = element("div");
    			div3 = element("div");
    			div2 = element("div");
    			div2.textContent = "ParallelCoin";
    			t3 = space();
    			div5 = element("div");
    			h1 = element("h1");
    			h1.textContent = "Svelte Router";
    			t5 = space();
    			div6 = element("div");
    			t6 = space();
    			div7 = element("div");
    			t7 = space();
    			div8 = element("div");
    			t8 = space();
    			div9 = element("div");
    			button = element("button");
    			button.textContent = "Open";
    			attr_dev(div0, "class", "h1");
    			add_location(div0, file$b, 15, 2, 189);
    			attr_dev(div1, "class", "h2");
    			add_location(div1, file$b, 16, 2, 214);
    			attr_dev(div2, "class", "analysis");
    			add_location(div2, file$b, 19, 4, 291);
    			attr_dev(div3, "class", "searchContent");
    			add_location(div3, file$b, 18, 3, 259);
    			attr_dev(div4, "class", "h3");
    			add_location(div4, file$b, 17, 2, 239);
    			add_location(h1, file$b, 24, 4, 375);
    			attr_dev(div5, "class", "h4");
    			add_location(div5, file$b, 22, 2, 353);
    			attr_dev(div6, "class", "h5");
    			add_location(div6, file$b, 31, 2, 424);
    			attr_dev(div7, "class", "h6");
    			add_location(div7, file$b, 32, 2, 449);
    			attr_dev(div8, "class", "h7");
    			add_location(div8, file$b, 33, 2, 474);
    			attr_dev(button, "id", "toggle");
    			attr_dev(button, "ref", "toggleBoardbtn");
    			attr_dev(button, "class", "e-btn e-info svelte-1bt0ocf");
    			attr_dev(button, "cssclass", "e-flat");
    			attr_dev(button, "iconcss", "e-icons burg-icon");
    			attr_dev(button, "istoggle", "true");
    			add_location(button, file$b, 35, 3, 519);
    			attr_dev(div9, "class", "h8");
    			add_location(div9, file$b, 34, 2, 499);
    			attr_dev(header, "class", "Header bgLight");
    			add_location(header, file$b, 14, 0, 155);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, header, anchor);
    			append_dev(header, div0);
    			append_dev(header, t0);
    			append_dev(header, div1);
    			append_dev(header, t1);
    			append_dev(header, div4);
    			append_dev(div4, div3);
    			append_dev(div3, div2);
    			append_dev(header, t3);
    			append_dev(header, div5);
    			append_dev(div5, h1);
    			append_dev(header, t5);
    			append_dev(header, div6);
    			append_dev(header, t6);
    			append_dev(header, div7);
    			append_dev(header, t7);
    			append_dev(header, div8);
    			append_dev(header, t8);
    			append_dev(header, div9);
    			append_dev(div9, button);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(header);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$h.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Header extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$h, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Header",
    			options,
    			id: create_fragment$h.name
    		});
    	}
    }

    /* src/components/ico/IcoOverview.svelte generated by Svelte v3.14.1 */

    const file$c = "src/components/ico/IcoOverview.svelte";

    function create_fragment$i(ctx) {
    	let svg;
    	let polygon0;
    	let g;
    	let polygon1;
    	let rect0;
    	let polygon2;
    	let rect1;
    	let rect2;
    	let path;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			polygon0 = svg_element("polygon");
    			g = svg_element("g");
    			polygon1 = svg_element("polygon");
    			rect0 = svg_element("rect");
    			polygon2 = svg_element("polygon");
    			rect1 = svg_element("rect");
    			rect2 = svg_element("rect");
    			path = svg_element("path");
    			attr_dev(polygon0, "class", "bgMoreLight");
    			attr_dev(polygon0, "points", "42,39 6,39 6,23 24,6 42,23");
    			add_location(polygon0, file$c, 0, 120, 120);
    			attr_dev(polygon1, "points", "39,21 34,16 34,9 39,9");
    			add_location(polygon1, file$c, 0, 204, 204);
    			attr_dev(rect0, "x", "6");
    			attr_dev(rect0, "y", "39");
    			attr_dev(rect0, "width", "36");
    			attr_dev(rect0, "height", "5");
    			add_location(rect0, file$c, 0, 245, 245);
    			attr_dev(g, "class", "bgBlue");
    			add_location(g, file$c, 0, 186, 186);
    			attr_dev(polygon2, "class", "red");
    			attr_dev(polygon2, "points", "24,4.3 4,22.9 6,25.1 24,8.4 42,25.1 44,22.9");
    			add_location(polygon2, file$c, 0, 291, 291);
    			attr_dev(rect1, "x", "18");
    			attr_dev(rect1, "y", "28");
    			attr_dev(rect1, "class", "red");
    			attr_dev(rect1, "width", "12");
    			attr_dev(rect1, "height", "16");
    			add_location(rect1, file$c, 0, 366, 366);
    			attr_dev(rect2, "x", "21");
    			attr_dev(rect2, "y", "17");
    			attr_dev(rect2, "class", "bgBlue");
    			attr_dev(rect2, "width", "6");
    			attr_dev(rect2, "height", "6");
    			add_location(rect2, file$c, 0, 422, 422);
    			attr_dev(path, "class", "bgGreen");
    			attr_dev(path, "d", "M27.5,35.5c-0.3,0-0.5,0.2-0.5,0.5v2c0,0.3,0.2,0.5,0.5,0.5S28,38.3,28,38v-2C28,35.7,27.8,35.5,27.5,35.5z");
    			add_location(path, file$c, 0, 479, 479);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$c, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, polygon0);
    			append_dev(svg, g);
    			append_dev(g, polygon1);
    			append_dev(g, rect0);
    			append_dev(svg, polygon2);
    			append_dev(svg, rect1);
    			append_dev(svg, rect2);
    			append_dev(svg, path);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$i.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoOverview extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$i, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoOverview",
    			options,
    			id: create_fragment$i.name
    		});
    	}
    }

    /* src/components/ico/IcoHistory.svelte generated by Svelte v3.14.1 */

    const file$d = "src/components/ico/IcoHistory.svelte";

    function create_fragment$j(ctx) {
    	let svg;
    	let path0;
    	let path1;
    	let g0;
    	let circle0;
    	let circle1;
    	let g1;
    	let path2;
    	let path3;
    	let g2;
    	let rect0;
    	let rect1;
    	let rect2;
    	let rect3;
    	let rect4;
    	let rect5;
    	let rect6;
    	let rect7;
    	let rect8;
    	let rect9;
    	let rect10;
    	let rect11;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			g0 = svg_element("g");
    			circle0 = svg_element("circle");
    			circle1 = svg_element("circle");
    			g1 = svg_element("g");
    			path2 = svg_element("path");
    			path3 = svg_element("path");
    			g2 = svg_element("g");
    			rect0 = svg_element("rect");
    			rect1 = svg_element("rect");
    			rect2 = svg_element("rect");
    			rect3 = svg_element("rect");
    			rect4 = svg_element("rect");
    			rect5 = svg_element("rect");
    			rect6 = svg_element("rect");
    			rect7 = svg_element("rect");
    			rect8 = svg_element("rect");
    			rect9 = svg_element("rect");
    			rect10 = svg_element("rect");
    			rect11 = svg_element("rect");
    			attr_dev(path0, "class", "bgMoreLight");
    			attr_dev(path0, "d", "M5,38V14h38v24c0,2.2-1.8,4-4,4H9C6.8,42,5,40.2,5,38z");
    			add_location(path0, file$d, 0, 119, 119);
    			attr_dev(path1, "class", "red");
    			attr_dev(path1, "d", "M43,10v6H5v-6c0-2.2,1.8-4,4-4h30C41.2,6,43,7.8,43,10z");
    			add_location(path1, file$d, 0, 203, 203);
    			attr_dev(circle0, "cx", "33");
    			attr_dev(circle0, "cy", "10");
    			attr_dev(circle0, "r", "3");
    			add_location(circle0, file$d, 0, 303, 303);
    			attr_dev(circle1, "cx", "15");
    			attr_dev(circle1, "cy", "10");
    			attr_dev(circle1, "r", "3");
    			add_location(circle1, file$d, 0, 334, 334);
    			attr_dev(g0, "class", "bgMoreLight");
    			add_location(g0, file$d, 0, 280, 280);
    			attr_dev(path2, "d", "M33,3c-1.1,0-2,0.9-2,2v5c0,1.1,0.9,2,2,2s2-0.9,2-2V5C35,3.9,34.1,3,33,3z");
    			add_location(path2, file$d, 0, 387, 387);
    			attr_dev(path3, "d", "M15,3c-1.1,0-2,0.9-2,2v5c0,1.1,0.9,2,2,2s2-0.9,2-2V5C17,3.9,16.1,3,15,3z");
    			add_location(path3, file$d, 0, 471, 471);
    			attr_dev(g1, "class", "bgGray");
    			add_location(g1, file$d, 0, 369, 369);
    			attr_dev(rect0, "x", "13");
    			attr_dev(rect0, "y", "20");
    			attr_dev(rect0, "width", "4");
    			attr_dev(rect0, "height", "4");
    			add_location(rect0, file$d, 0, 577, 577);
    			attr_dev(rect1, "x", "19");
    			attr_dev(rect1, "y", "20");
    			attr_dev(rect1, "width", "4");
    			attr_dev(rect1, "height", "4");
    			add_location(rect1, file$d, 0, 619, 619);
    			attr_dev(rect2, "x", "25");
    			attr_dev(rect2, "y", "20");
    			attr_dev(rect2, "width", "4");
    			attr_dev(rect2, "height", "4");
    			add_location(rect2, file$d, 0, 661, 661);
    			attr_dev(rect3, "x", "31");
    			attr_dev(rect3, "y", "20");
    			attr_dev(rect3, "width", "4");
    			attr_dev(rect3, "height", "4");
    			add_location(rect3, file$d, 0, 703, 703);
    			attr_dev(rect4, "x", "13");
    			attr_dev(rect4, "y", "26");
    			attr_dev(rect4, "width", "4");
    			attr_dev(rect4, "height", "4");
    			add_location(rect4, file$d, 0, 745, 745);
    			attr_dev(rect5, "x", "19");
    			attr_dev(rect5, "y", "26");
    			attr_dev(rect5, "width", "4");
    			attr_dev(rect5, "height", "4");
    			add_location(rect5, file$d, 0, 787, 787);
    			attr_dev(rect6, "x", "25");
    			attr_dev(rect6, "y", "26");
    			attr_dev(rect6, "width", "4");
    			attr_dev(rect6, "height", "4");
    			add_location(rect6, file$d, 0, 829, 829);
    			attr_dev(rect7, "x", "31");
    			attr_dev(rect7, "y", "26");
    			attr_dev(rect7, "width", "4");
    			attr_dev(rect7, "height", "4");
    			add_location(rect7, file$d, 0, 871, 871);
    			attr_dev(rect8, "x", "13");
    			attr_dev(rect8, "y", "32");
    			attr_dev(rect8, "width", "4");
    			attr_dev(rect8, "height", "4");
    			add_location(rect8, file$d, 0, 913, 913);
    			attr_dev(rect9, "x", "19");
    			attr_dev(rect9, "y", "32");
    			attr_dev(rect9, "width", "4");
    			attr_dev(rect9, "height", "4");
    			add_location(rect9, file$d, 0, 955, 955);
    			attr_dev(rect10, "x", "25");
    			attr_dev(rect10, "y", "32");
    			attr_dev(rect10, "width", "4");
    			attr_dev(rect10, "height", "4");
    			add_location(rect10, file$d, 0, 997, 997);
    			attr_dev(rect11, "x", "31");
    			attr_dev(rect11, "y", "32");
    			attr_dev(rect11, "width", "4");
    			attr_dev(rect11, "height", "4");
    			add_location(rect11, file$d, 0, 1039, 1039);
    			attr_dev(g2, "class", "bgGray");
    			add_location(g2, file$d, 0, 559, 559);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$d, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path0);
    			append_dev(svg, path1);
    			append_dev(svg, g0);
    			append_dev(g0, circle0);
    			append_dev(g0, circle1);
    			append_dev(svg, g1);
    			append_dev(g1, path2);
    			append_dev(g1, path3);
    			append_dev(svg, g2);
    			append_dev(g2, rect0);
    			append_dev(g2, rect1);
    			append_dev(g2, rect2);
    			append_dev(g2, rect3);
    			append_dev(g2, rect4);
    			append_dev(g2, rect5);
    			append_dev(g2, rect6);
    			append_dev(g2, rect7);
    			append_dev(g2, rect8);
    			append_dev(g2, rect9);
    			append_dev(g2, rect10);
    			append_dev(g2, rect11);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$j.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoHistory extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$j, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoHistory",
    			options,
    			id: create_fragment$j.name
    		});
    	}
    }

    /* src/components/ico/IcoAddressBook.svelte generated by Svelte v3.14.1 */

    const file$e = "src/components/ico/IcoAddressBook.svelte";

    function create_fragment$k(ctx) {
    	let svg;
    	let path0;
    	let path1;
    	let path2;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			path2 = svg_element("path");
    			attr_dev(path0, "class", "bgLightBlue");
    			attr_dev(path0, "d", "M38,44H12V4h26c2.2,0,4,1.8,4,4v32C42,42.2,40.2,44,38,44z");
    			add_location(path0, file$e, 0, 119, 119);
    			attr_dev(path1, "class", "bgBlue");
    			attr_dev(path1, "d", "M10,4h2v40h-2c-2.2,0-4-1.8-4-4V8C6,5.8,7.8,4,10,4z");
    			add_location(path1, file$e, 0, 207, 207);
    			attr_dev(path2, "class", "bgMoreLight");
    			attr_dev(path2, "d", "M36,24.2c-0.1,4.8-3.1,6.9-5.3,6.7c-0.6-0.1-2.1-0.1-2.9-1.6c-0.8,1-1.8,1.6-3.1,1.6c-2.6,0-3.3-2.5-3.4-3.1 c-0.1-0.7-0.2-1.4-0.1-2.2c0.1-1,1.1-6.5,5.7-6.5c2.2,0,3.5,1.1,3.7,1.3L30,27.2c0,0.3-0.2,1.6,1.1,1.6c2.1,0,2.4-3.9,2.4-4.6 c0.1-1.2,0.3-8.2-7-8.2c-6.9,0-7.9,7.4-8,9.2c-0.5,8.5,6,8.5,7.2,8.5c1.7,0,3.7-0.7,3.9-0.8l0.4,2c-0.3,0.2-2,1.1-4.4,1.1 c-2.2,0-10.1-0.4-9.8-10.8C16.1,23.1,17.4,14,26.6,14C35.8,14,36,22.1,36,24.2z M24.1,25.5c-0.1,1,0,1.8,0.2,2.3 c0.2,0.5,0.6,0.8,1.2,0.8c0.1,0,0.3,0,0.4-0.1c0.2-0.1,0.3-0.1,0.5-0.3c0.2-0.1,0.3-0.3,0.5-0.6c0.2-0.2,0.3-0.6,0.4-1l0.5-5.4 c-0.2-0.1-0.5-0.1-0.7-0.1c-0.5,0-0.9,0.1-1.2,0.3c-0.3,0.2-0.6,0.5-0.9,0.8c-0.2,0.4-0.4,0.8-0.6,1.3S24.2,24.8,24.1,25.5z");
    			add_location(path2, file$e, 0, 284, 284);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$e, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path0);
    			append_dev(svg, path1);
    			append_dev(svg, path2);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$k.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoAddressBook extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$k, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoAddressBook",
    			options,
    			id: create_fragment$k.name
    		});
    	}
    }

    /* src/components/ico/IcoSettings.svelte generated by Svelte v3.14.1 */

    const file$f = "src/components/ico/IcoSettings.svelte";

    function create_fragment$l(ctx) {
    	let svg;
    	let path0;
    	let path1;

    	const block = {
    		c: function create() {
    			svg = svg_element("svg");
    			path0 = svg_element("path");
    			path1 = svg_element("path");
    			attr_dev(path0, "class", "bgGray");
    			attr_dev(path0, "d", "M39.6,27.2c0.1-0.7,0.2-1.4,0.2-2.2s-0.1-1.5-0.2-2.2l4.5-3.2c0.4-0.3,0.6-0.9,0.3-1.4L40,10.8 c-0.3-0.5-0.8-0.7-1.3-0.4l-5,2.3c-1.2-0.9-2.4-1.6-3.8-2.2l-0.5-5.5c-0.1-0.5-0.5-0.9-1-0.9h-8.6c-0.5,0-1,0.4-1,0.9l-0.5,5.5 c-1.4,0.6-2.7,1.3-3.8,2.2l-5-2.3c-0.5-0.2-1.1,0-1.3,0.4l-4.3,7.4c-0.3,0.5-0.1,1.1,0.3,1.4l4.5,3.2c-0.1,0.7-0.2,1.4-0.2,2.2 s0.1,1.5,0.2,2.2L4,30.4c-0.4,0.3-0.6,0.9-0.3,1.4L8,39.2c0.3,0.5,0.8,0.7,1.3,0.4l5-2.3c1.2,0.9,2.4,1.6,3.8,2.2l0.5,5.5 c0.1,0.5,0.5,0.9,1,0.9h8.6c0.5,0,1-0.4,1-0.9l0.5-5.5c1.4-0.6,2.7-1.3,3.8-2.2l5,2.3c0.5,0.2,1.1,0,1.3-0.4l4.3-7.4 c0.3-0.5,0.1-1.1-0.3-1.4L39.6,27.2z M24,35c-5.5,0-10-4.5-10-10c0-5.5,4.5-10,10-10c5.5,0,10,4.5,10,10C34,30.5,29.5,35,24,35z");
    			add_location(path0, file$f, 0, 120, 120);
    			attr_dev(path1, "class", "black");
    			attr_dev(path1, "d", "M24,13c-6.6,0-12,5.4-12,12c0,6.6,5.4,12,12,12s12-5.4,12-12C36,18.4,30.6,13,24,13z M24,30 c-2.8,0-5-2.2-5-5c0-2.8,2.2-5,5-5s5,2.2,5,5C29,27.8,26.8,30,24,30z");
    			add_location(path1, file$f, 0, 839, 839);
    			attr_dev(svg, "version", "1");
    			attr_dev(svg, "xmlns", "http://www.w3.org/2000/svg");
    			attr_dev(svg, "class", "icon");
    			attr_dev(svg, "viewBox", "0 0 48 48");
    			attr_dev(svg, "enable-background", "new 0 0 48 48");
    			add_location(svg, file$f, 0, 0, 0);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, svg, anchor);
    			append_dev(svg, path0);
    			append_dev(svg, path1);
    		},
    		p: noop,
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(svg);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$l.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class IcoSettings extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$l, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "IcoSettings",
    			options,
    			id: create_fragment$l.name
    		});
    	}
    }

    /* src/components/layout/Nav.svelte generated by Svelte v3.14.1 */
    const file$g = "src/components/layout/Nav.svelte";

    function create_fragment$m(ctx) {
    	let nav;
    	let ul;
    	let li0;
    	let button0;
    	let t0;
    	let li1;
    	let button1;
    	let t1;
    	let li2;
    	let button2;
    	let t2;
    	let li3;
    	let button3;
    	let t3;
    	let li4;
    	let button4;
    	let current;
    	let dispose;
    	const icooverview0 = new IcoOverview({ $$inline: true });
    	const icohistory = new IcoHistory({ $$inline: true });
    	const icoaddressbook = new IcoAddressBook({ $$inline: true });
    	const icooverview1 = new IcoOverview({ $$inline: true });
    	const icosettings = new IcoSettings({ $$inline: true });

    	const block = {
    		c: function create() {
    			nav = element("nav");
    			ul = element("ul");
    			li0 = element("li");
    			button0 = element("button");
    			create_component(icooverview0.$$.fragment);
    			t0 = space();
    			li1 = element("li");
    			button1 = element("button");
    			create_component(icohistory.$$.fragment);
    			t1 = space();
    			li2 = element("li");
    			button2 = element("button");
    			create_component(icoaddressbook.$$.fragment);
    			t2 = space();
    			li3 = element("li");
    			button3 = element("button");
    			create_component(icooverview1.$$.fragment);
    			t3 = space();
    			li4 = element("li");
    			button4 = element("button");
    			create_component(icosettings.$$.fragment);
    			attr_dev(button0, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button0, file$g, 14, 10, 418);
    			attr_dev(li0, "id", "menuoverview");
    			attr_dev(li0, "class", "sidebar-item current");
    			add_location(li0, file$g, 13, 6, 356);
    			attr_dev(button1, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button1, file$g, 19, 8, 637);
    			attr_dev(li1, "id", "menutransactions");
    			attr_dev(li1, "class", "sidebar-item");
    			add_location(li1, file$g, 18, 6, 581);
    			attr_dev(button2, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button2, file$g, 24, 8, 856);
    			attr_dev(li2, "id", "menuaddressbook");
    			attr_dev(li2, "class", "sidebar-item");
    			add_location(li2, file$g, 23, 6, 801);
    			attr_dev(button3, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button3, file$g, 29, 8, 1082);
    			attr_dev(li3, "id", "menublockexplorer");
    			attr_dev(li3, "class", "sidebar-item");
    			add_location(li3, file$g, 28, 6, 1025);
    			attr_dev(button4, "class", "noMargin noPadding noBorder bgTrans sXs cursorPointer");
    			add_location(button4, file$g, 34, 8, 1293);
    			attr_dev(li4, "id", "menusettings");
    			attr_dev(li4, "class", "sidebar-item");
    			add_location(li4, file$g, 33, 6, 1241);
    			attr_dev(ul, "id", "menu");
    			attr_dev(ul, "class", "lsn noPadding");
    			add_location(ul, file$g, 12, 4, 313);
    			attr_dev(nav, "class", "Nav");
    			add_location(nav, file$g, 11, 0, 291);

    			dispose = [
    				listen_dev(button0, "click", isPage.overview, false, false, false),
    				listen_dev(button1, "click", isPage.transactions, false, false, false),
    				listen_dev(button2, "click", isPage.addressbook, false, false, false),
    				listen_dev(button3, "click", isPage.explorer, false, false, false),
    				listen_dev(button4, "click", isPage.settings, false, false, false)
    			];
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, nav, anchor);
    			append_dev(nav, ul);
    			append_dev(ul, li0);
    			append_dev(li0, button0);
    			mount_component(icooverview0, button0, null);
    			append_dev(ul, t0);
    			append_dev(ul, li1);
    			append_dev(li1, button1);
    			mount_component(icohistory, button1, null);
    			append_dev(ul, t1);
    			append_dev(ul, li2);
    			append_dev(li2, button2);
    			mount_component(icoaddressbook, button2, null);
    			append_dev(ul, t2);
    			append_dev(ul, li3);
    			append_dev(li3, button3);
    			mount_component(icooverview1, button3, null);
    			append_dev(ul, t3);
    			append_dev(ul, li4);
    			append_dev(li4, button4);
    			mount_component(icosettings, button4, null);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(icooverview0.$$.fragment, local);
    			transition_in(icohistory.$$.fragment, local);
    			transition_in(icoaddressbook.$$.fragment, local);
    			transition_in(icooverview1.$$.fragment, local);
    			transition_in(icosettings.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(icooverview0.$$.fragment, local);
    			transition_out(icohistory.$$.fragment, local);
    			transition_out(icoaddressbook.$$.fragment, local);
    			transition_out(icooverview1.$$.fragment, local);
    			transition_out(icosettings.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(nav);
    			destroy_component(icooverview0);
    			destroy_component(icohistory);
    			destroy_component(icoaddressbook);
    			destroy_component(icooverview1);
    			destroy_component(icosettings);
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$m.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class Nav extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$m, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "Nav",
    			options,
    			id: create_fragment$m.name
    		});
    	}
    }

    /* src/components/layout/UI.svelte generated by Svelte v3.14.1 */
    const file$h = "src/components/layout/UI.svelte";

    function create_fragment$n(ctx) {
    	let t0;
    	let div7;
    	let div6;
    	let div5;
    	let div0;
    	let t1;
    	let t2;
    	let div3;
    	let div1;
    	let t3;
    	let t4;
    	let div2;
    	let t5;
    	let div4;
    	let div4_transition;
    	let current;
    	const logo = new Logo({ $$inline: true });
    	const header = new Header({ $$inline: true });
    	const nav = new Nav({ $$inline: true });
    	var switch_value = ctx.$isPage;

    	function switch_props(ctx) {
    		return { $$inline: true };
    	}

    	if (switch_value) {
    		var switch_instance = new switch_value(switch_props());
    	}

    	const block = {
    		c: function create() {
    			t0 = space();
    			div7 = element("div");
    			div6 = element("div");
    			div5 = element("div");
    			div0 = element("div");
    			create_component(logo.$$.fragment);
    			t1 = space();
    			create_component(header.$$.fragment);
    			t2 = space();
    			div3 = element("div");
    			div1 = element("div");
    			t3 = space();
    			create_component(nav.$$.fragment);
    			t4 = space();
    			div2 = element("div");
    			t5 = space();
    			div4 = element("div");
    			if (switch_instance) create_component(switch_instance.$$.fragment);
    			document.title = "DuOS";
    			attr_dev(div0, "class", "flx fii Logo");
    			add_location(div0, file$h, 30, 5, 513);
    			attr_dev(div1, "class", "Open");
    			add_location(div1, file$h, 33, 6, 611);
    			attr_dev(div2, "class", "Side");
    			add_location(div2, file$h, 35, 6, 656);
    			attr_dev(div3, "class", "Sidebar bgLight");
    			add_location(div3, file$h, 32, 5, 575);
    			attr_dev(div4, "id", "main");
    			attr_dev(div4, "class", "grayGrad Main");
    			add_location(div4, file$h, 38, 0, 694);
    			attr_dev(div5, "class", "grid-container rwrap bgDark");
    			add_location(div5, file$h, 29, 4, 466);
    			attr_dev(div6, "id", "display");
    			attr_dev(div6, "class", "fii");
    			add_location(div6, file$h, 28, 53, 431);
    			attr_dev(div7, "id", "x");
    			attr_dev(div7, "class", "fullScreen bgDark flx lightTheme");
    			add_location(div7, file$h, 28, 0, 378);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, t0, anchor);
    			insert_dev(target, div7, anchor);
    			append_dev(div7, div6);
    			append_dev(div6, div5);
    			append_dev(div5, div0);
    			mount_component(logo, div0, null);
    			append_dev(div5, t1);
    			mount_component(header, div5, null);
    			append_dev(div5, t2);
    			append_dev(div5, div3);
    			append_dev(div3, div1);
    			append_dev(div3, t3);
    			mount_component(nav, div3, null);
    			append_dev(div3, t4);
    			append_dev(div3, div2);
    			append_dev(div5, t5);
    			append_dev(div5, div4);

    			if (switch_instance) {
    				mount_component(switch_instance, div4, null);
    			}

    			current = true;
    		},
    		p: function update(changed, ctx) {
    			if (switch_value !== (switch_value = ctx.$isPage)) {
    				if (switch_instance) {
    					group_outros();
    					const old_component = switch_instance;

    					transition_out(old_component.$$.fragment, 1, 0, () => {
    						destroy_component(old_component, 1);
    					});

    					check_outros();
    				}

    				if (switch_value) {
    					switch_instance = new switch_value(switch_props());
    					create_component(switch_instance.$$.fragment);
    					transition_in(switch_instance.$$.fragment, 1);
    					mount_component(switch_instance, div4, null);
    				} else {
    					switch_instance = null;
    				}
    			}
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(logo.$$.fragment, local);
    			transition_in(header.$$.fragment, local);
    			transition_in(nav.$$.fragment, local);
    			if (switch_instance) transition_in(switch_instance.$$.fragment, local);

    			add_render_callback(() => {
    				if (!div4_transition) div4_transition = create_bidirectional_transition(div4, fade, { duration: 300 }, true);
    				div4_transition.run(1);
    			});

    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(logo.$$.fragment, local);
    			transition_out(header.$$.fragment, local);
    			transition_out(nav.$$.fragment, local);
    			if (switch_instance) transition_out(switch_instance.$$.fragment, local);
    			if (!div4_transition) div4_transition = create_bidirectional_transition(div4, fade, { duration: 300 }, false);
    			div4_transition.run(0);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(t0);
    			if (detaching) detach_dev(div7);
    			destroy_component(logo);
    			destroy_component(header);
    			destroy_component(nav);
    			if (switch_instance) destroy_component(switch_instance);
    			if (detaching && div4_transition) div4_transition.end();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$n.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$5($$self, $$props, $$invalidate) {
    	let $isPage;
    	validate_store(isPage, "isPage");
    	component_subscribe($$self, isPage, $$value => $$invalidate("$isPage", $isPage = $$value));
    	let page;

    	const unsubscribe = isPage.subscribe(value => {
    		page = value;
    	});

    	$$self.$capture_state = () => {
    		return {};
    	};

    	$$self.$inject_state = $$props => {
    		if ("page" in $$props) page = $$props.page;
    		if ("$isPage" in $$props) isPage.set($isPage = $$props.$isPage);
    	};

    	return { $isPage };
    }

    class UI extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$5, create_fragment$n, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "UI",
    			options,
    			id: create_fragment$n.name
    		});
    	}
    }

    /* src/DuOS.svelte generated by Svelte v3.14.1 */

    // (9:0) {#if bios.boot}
    function create_if_block_1(ctx) {
    	let current;
    	const boot = new Boot({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(boot.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(boot, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(boot.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(boot.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(boot, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_1.name,
    		type: "if",
    		source: "(9:0) {#if bios.boot}",
    		ctx
    	});

    	return block;
    }

    // (12:0) {#if !bios.boot}
    function create_if_block$1(ctx) {
    	let current;
    	const ui = new UI({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(ui.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(ui, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(ui.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(ui.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(ui, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block$1.name,
    		type: "if",
    		source: "(12:0) {#if !bios.boot}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$o(ctx) {
    	let t;
    	let if_block1_anchor;
    	let current;
    	let if_block0 = bios.boot && create_if_block_1(ctx);
    	let if_block1 = !bios.boot && create_if_block$1(ctx);

    	const block = {
    		c: function create() {
    			if (if_block0) if_block0.c();
    			t = space();
    			if (if_block1) if_block1.c();
    			if_block1_anchor = empty();
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			if (if_block0) if_block0.m(target, anchor);
    			insert_dev(target, t, anchor);
    			if (if_block1) if_block1.m(target, anchor);
    			insert_dev(target, if_block1_anchor, anchor);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block0);
    			transition_in(if_block1);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block0);
    			transition_out(if_block1);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (if_block0) if_block0.d(detaching);
    			if (detaching) detach_dev(t);
    			if (if_block1) if_block1.d(detaching);
    			if (detaching) detach_dev(if_block1_anchor);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$o.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    class DuOS extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, null, create_fragment$o, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "DuOS",
    			options,
    			id: create_fragment$o.name
    		});
    	}
    }

    const duos = new DuOS({
    	target: document.body
    });

    return duos;

}(smelte));
//# sourceMappingURL=dui.js.map
